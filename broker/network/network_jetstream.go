// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package network

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"

	"github.com/choria-io/go-choria/backoff"
	"github.com/choria-io/go-choria/scout"
)

func (s *Server) setupStreaming() error {
	if s.config.Choria.NetworkStreamStore == "" {
		return nil
	}

	if s.gnatsd.SystemAccount() == nil {
		return fmt.Errorf("system Account is required for Choria Streams")
	}

	s.log.Infof("Enabling Choria Streams for account %s", s.choriaAccount)

	err := s.choriaAccount.EnableJetStream(nil)
	if err != nil {
		s.log.Errorf("Could not enable Choria Streams for the %s account: %s", s.choriaAccount.Name, err)
	}

	if !s.choriaAccount.JetStreamEnabled() {
		s.log.Errorf("Choria Streams enabled for account %q but it's not reporting as enabled", s.choriaAccount.Name)
	}

	return nil
}

func (s *Server) configureSystemStreams(ctx context.Context) error {
	if s.config.Choria.NetworkStreamStore == "" {
		return nil
	}

	var opts []nats.Option

	if s.IsTLS() {
		s.log.Info("Configuring Choria System Streams with TLS")
		tlsc, err := s.choria.ClientTLSConfig()
		if err != nil {
			return err
		}
		opts = append(opts, nats.Secure(tlsc))
	} else {
		s.log.Info("Configuring Choria System Streams without TLS")
	}

	var nc *nats.Conn
	var err error

	err = backoff.TwentySec.For(ctx, func(try int) error {
		nc, err = nats.Connect(s.opts.ClientAdvertise, opts...)
		if err != nil {
			s.log.Warnf("Could not connect to broker %s to configure System Streams: %s", s.opts.ClientAdvertise, err)
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	defer nc.Close()

	mgr, err := jsm.New(nc)
	if err != nil {
		return err
	}

	cfg := s.config.Choria
	if cfg.NetworkEventStoreReplicas == -1 || cfg.NetworkMachineStoreReplicas == -1 || cfg.NetworkStreamAdvisoryReplicas == -1 || cfg.NetworkLeaderElectionReplicas == -1 {
		delay := time.Duration(rand.Intn(60)+10) * time.Second
		s.log.Infof("Configuring system streams after %v", delay)
		err = backoff.Default.Sleep(ctx, delay)
		if err != nil {
			s.log.Errorf("Aborting stream configuration: %v", err)
			return err
		}

		peers, err := s.choria.NetworkBrokerPeers()
		if err != nil {
			s.log.Warnf("Cannot determine network peers to calculate dynamic replica sizes: %s", err)
		}

		count := peers.Count()
		if count == 0 {
			count = 1 // avoid replica=0
		}

		if cfg.NetworkEventStoreReplicas == -1 {
			s.log.Infof("Setting Lifecycle Event Store Replicas to %d", count)
			cfg.NetworkEventStoreReplicas = count
		}

		if cfg.NetworkMachineStoreReplicas == -1 {
			s.log.Infof("Setting Autonomous Agent Event Store Replicas to %d", count)
			cfg.NetworkMachineStoreReplicas = count
		}

		if cfg.NetworkStreamAdvisoryReplicas == -1 {
			s.log.Infof("Setting Choria Streams Advisory Store Replicas to %d", count)
			cfg.NetworkStreamAdvisoryReplicas = count
		}

		if cfg.NetworkLeaderElectionReplicas == -1 {
			s.log.Infof("Setting Choria Streams Leader election Replicas to %d", count)
			cfg.NetworkLeaderElectionReplicas = count
		}
	}

	err = s.createOrUpdateStream("CHORIA_EVENTS", []string{"choria.lifecycle.>"}, cfg.NetworkEventStoreDuration, cfg.NetworkEventStoreReplicas, mgr)
	if err != nil {
		return err
	}

	err = s.createOrUpdateStream("CHORIA_MACHINE", []string{"choria.machine.>"}, cfg.NetworkMachineStoreDuration, cfg.NetworkMachineStoreReplicas, mgr)
	if err != nil {
		return err
	}

	err = s.createOrUpdateStream("CHORIA_STREAM_ADVISORIES", []string{"$JS.EVENT.ADVISORY.>"}, cfg.NetworkStreamAdvisoryDuration, cfg.NetworkStreamAdvisoryReplicas, mgr)
	if err != nil {
		return err
	}

	err = scout.ConfigureStreams(nc, s.log.WithField("component", "scout"))
	if err != nil {
		return err
	}

	eCfg, err := jsm.NewStreamConfiguration(jsm.DefaultStream,
		jsm.Replicas(cfg.NetworkLeaderElectionReplicas),
		jsm.MaxAge(cfg.NetworkLeaderElectionTTL),
		jsm.AllowRollup(),
		jsm.Subjects("$KV.CHORIA_LEADER_ELECTION.>"),
		jsm.StreamDescription("Choria Leader Election Bucket"),
		jsm.MaxMessageSize(1024),
		jsm.FileStorage(),
		jsm.MaxMessagesPerSubject(1))
	if err != nil {
		return err
	}
	err = s.createOrUpdateStreamWithConfig("KV_CHORIA_LEADER_ELECTION", *eCfg, mgr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) createOrUpdateStream(name string, subjects []string, maxAge time.Duration, replicas int, mgr *jsm.Manager) error {
	if int(maxAge) <= 0 {
		return nil
	}

	cfg, err := jsm.NewStreamConfiguration(jsm.DefaultStream, jsm.FileStorage(), jsm.Subjects(subjects...), jsm.MaxAge(maxAge), jsm.Replicas(replicas))
	if err != nil {
		return fmt.Errorf("could not create configuration: %s", err)
	}

	err = s.createOrUpdateStreamWithConfig(name, *cfg, mgr)
	if err != nil {
		return fmt.Errorf("could not create stream %s: %s", name, err)
	}

	return nil
}

func (s *Server) createOrUpdateStreamWithConfig(name string, cfg api.StreamConfig, mgr *jsm.Manager) error {
	cfg.Name = name
	str, err := mgr.LoadStream(name)
	if err != nil {
		_, err := mgr.NewStreamFromDefault(name, cfg)
		if err == nil {
			s.log.Infof("Created stream %s with %d replicas and %s retention", cfg.Name, cfg.Replicas, cfg.MaxAge)
		}
		return err
	}

	err = str.UpdateConfiguration(cfg)
	if err != nil {
		return err
	}

	s.log.Infof("Configured stream %s with %d replicas and %s retention", cfg.Name, cfg.Replicas, cfg.MaxAge)

	return nil
}
