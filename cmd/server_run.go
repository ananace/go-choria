// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/choria-io/go-choria/config"
	"github.com/choria-io/go-choria/internal/util"
	"github.com/choria-io/go-choria/protocol"
	"github.com/choria-io/go-choria/providers/provtarget"
	"github.com/choria-io/go-choria/server"
	log "github.com/sirupsen/logrus"
)

func (r *serverRunCommand) Setup() (err error) {
	if broker, ok := cmdWithFullCommand("server"); ok {
		r.cmd = broker.Cmd().Command("run", "Runs a Choria Server").Default()
		r.cmd.Flag("disable-tls", "Disables TLS").Hidden().Default("false").BoolVar(&r.disableTLS)
		r.cmd.Flag("disable-ssl-verification", "Disables SSL Verification").Hidden().Default("false").BoolVar(&r.disableTLSVerify)
		r.cmd.Flag("pid", "Write running PID to a file").StringVar(&r.pidFile)
		r.cmd.Flag("service-host", "Runs as a Service Agent host").BoolVar(&r.serviceHost)
	}

	return
}

func (r *serverRunCommand) Configure() error {
	if debug {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		log.Debug("Logging at debug level due to CLI override")
	}

	if configFile == "" {
		return fmt.Errorf("server run requires a configuration file")
	}

	switch {
	case util.FileExist(configFile):
		cfg, err = config.NewSystemConfig(configFile, true)
		if err != nil {
			return fmt.Errorf("could not parse configuration: %s", err)
		}

		provtarget.Configure(cfg, log.WithField("component", "provtarget"))

		// if a config file existed and prov is disable even after reading it, discard it and start fresh
		if r.shouldProvision(cfg) {
			log.Warnf("Switching to provisioning configuration due to build defaults and server.provision configuration setting")
			cfg, err = r.provisionConfig(configFile)
			if err != nil {
				return err
			}
		}

	case bi.ProvisionBrokerURLs() != "" || util.FileExist(bi.ProvisionJWTFile()):
		cfg, err = r.provisionConfig(configFile)
		if err != nil {
			return err
		}

	default:
		// we have no configuration file or anything, so we use defaults and possibly initiate provisioning
		cfg, err = config.NewDefaultSystemConfig(true)
		if err != nil {
			return fmt.Errorf("could not create default server configuration")
		}

		provtarget.Configure(cfg, log.WithField("component", "provtarget"))

		// if a config file didn't exist and prov is disabled we cant start
		if !r.shouldProvision(cfg) {
			return fmt.Errorf("configuration file %s was not found and provisioning is disabled", configFile)
		}

		log.Warnf("Switching to provisioning configuration due to build defaults and missing %s", configFile)

		cfg, err = r.provisionConfig(configFile)
		if err != nil {
			return err
		}
	}

	cfg.ApplyBuildSettings(bi)

	cfg.DisableSecurityProviderVerify = true

	if os.Getenv("INSECURE_YES_REALLY") == "true" {
		protocol.Secure = "false"
		cfg.DisableTLS = true
	}

	return nil
}

func (r *serverRunCommand) shouldProvision(cfg *config.Config) bool {
	prov := bi.ProvisionDefault()
	hasOpt := cfg.HasOption("plugin.choria.server.provision")
	if hasOpt {
		prov = cfg.Choria.Provision
	}
	return prov
}

func (r *serverRunCommand) provisionConfig(f string) (*config.Config, error) {
	cfg, err = config.NewDefaultSystemConfig(true)
	if err != nil {
		return nil, fmt.Errorf("could not create default configuration for provisioning: %s", err)
	}
	cfg.ConfigFile = f

	// set this to avoid calling into puppet on non puppet machines
	// later ConfigureProvisioning() will do all the right things
	cfg.Choria.SecurityProvider = "file"

	return cfg, nil
}

func (r *serverRunCommand) prepareInstance() (i *server.Instance, err error) {
	if r.disableTLS {
		c.Config.DisableTLS = true
		log.Warn("Running with TLS disabled, not compatible with production use.")
	}

	if r.disableTLSVerify {
		c.Config.DisableTLSVerify = true
		log.Warn("Running with TLS Verification disabled, not compatible with production use.")
	}

	c.ConfigureProvisioning()

	instance, err := server.NewInstance(c)
	if err != nil {
		return nil, fmt.Errorf("could not create Choria Server instance: %s", err)
	}

	log.Infof("Choria Server version %s starting with config %s", bi.Version(), c.Config.ConfigFile)

	if r.pidFile != "" {
		err := os.WriteFile(r.pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
		if err != nil {
			return nil, fmt.Errorf("could not write PID: %s", err)
		}
	}

	return instance, nil
}

func init() {
	cli.commands = append(cli.commands, &serverRunCommand{})
}
