// Copyright (c) 2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package streams

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/choria-io/go-choria/broker/adapter/ingest"
	"github.com/choria-io/go-choria/broker/adapter/stats"
	"github.com/choria-io/go-choria/inter"
	"github.com/choria-io/go-choria/protocol"
	"github.com/choria-io/go-choria/srvcache"

	"github.com/choria-io/go-choria/config"
	log "github.com/sirupsen/logrus"
)

// Streams is an adapter that connects a NATS topic with messages sent from Choria
// in its usual transport protocol to a Choria Streams stream.
//
// On the stream the messages will be JSON format with keys
// body, sender and time.  Body is a base64 encoded string
//
// Configure the adapters:
//   # required
//   plugin.choria.adapters = discovery
//   plugin.choria.adapter.discovery.type = choria_streams
//   plugin.choria.adapter.discovery.queue_len = 1000 # default
//
// Configure the stream output:
//
//   plugin.choria.adapter.discovery.stream.servers = js1:4222,js2:4222 # uses normal middleware server resolution when unset
//   plugin.choria.adapter.discovery.stream.topic = discovery # default, %s gets replaced with sender id
//   plugin.choria.adapter.discovery.stream.workers = 10 # default
//
// Configure the NATS ingest:
//
//    plugin.choria.adapter.discovery.ingest.topic = mcollective.broadcast.agent.discovery
//    plugin.choria.adapter.discovery.ingest.protocol = request # or reply
//    plugin.choria.adapter.discovery.ingest.workers = 10 # default
type Streams struct {
	streams []*stream
	ingests []*ingest.NatsIngest
	work    chan ingest.Adaptable
	log     *log.Entry
}

type Framework interface {
	Configuration() *config.Config
	MiddlewareServers() (servers srvcache.Servers, err error)
	NewConnector(ctx context.Context, servers func() (srvcache.Servers, error), name string, logger *log.Entry) (conn inter.Connector, err error)
	NewRequestFromTransportJSON(payload []byte, skipvalidate bool) (msg protocol.Request, err error)
	NewReplyFromTransportJSON(payload []byte, skipvalidate bool) (msg protocol.Reply, err error)
}

var fw Framework
var cfg *config.Config

func Create(name string, choria Framework) (adapter *Streams, err error) {
	fw = choria
	cfg = fw.Configuration()

	s := fmt.Sprintf("plugin.choria.adapter.%s.queue_len", name)
	worklen, err := strconv.Atoi(cfg.Option(s, "1000"))
	if err != nil {
		return nil, fmt.Errorf("%s should be a integer number", s)
	}

	stats.WorkQueueCapacityGauge.WithLabelValues(name, cfg.Identity).Set(float64(worklen))

	adapter = &Streams{
		log:  log.WithFields(log.Fields{"component": "streams_adapter", "name": name}),
		work: make(chan ingest.Adaptable, worklen),
	}

	adapter.ingests, err = ingest.New(name, adapter.work, choria, adapter.log)
	if err != nil {
		return nil, fmt.Errorf("could not create adapter %s: %s", name, err)
	}

	adapter.streams, err = newStream(name, adapter.work, adapter.log)
	if err != nil {
		return nil, fmt.Errorf("could not create adapter %s: %s", name, err)
	}

	return adapter, nil
}

func (sa *Streams) Init(ctx context.Context, cm inter.ConnectionManager) (err error) {
	for _, worker := range sa.streams {
		if ctx.Err() != nil {
			return fmt.Errorf("shutdown called")
		}

		err = worker.connect(ctx, cm)
		if err != nil {
			return fmt.Errorf("failure during initial Choria Streams connections: %s", err)
		}
	}

	for _, worker := range sa.ingests {
		if ctx.Err() != nil {
			return fmt.Errorf("shutdown called")
		}

		err = worker.Connect(ctx, cm)
		if err != nil {
			return fmt.Errorf("failure during Choria Streams initial connections: %s", err)
		}
	}

	return nil
}

func (sa *Streams) Process(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, worker := range sa.streams {
		wg.Add(1)
		go worker.publisher(ctx, wg)
	}

	for _, worker := range sa.ingests {
		wg.Add(1)
		go worker.Receiver(ctx, wg)
	}
}
