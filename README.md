# Choria Broker and Server

This is a daemon written in Go that will replace the traditional combination of Ruby MCollectived + Choria plugins - at least as far as the daemon side goes.

It will include at least:

  * A vendored `NATS` instance fully managed by `server.cfg` with clustering support
  * A Federation Broker for the Choria protocol
  * An Protocol Adapter framework to rewrite message into other systems.  For eg. Discovery messages to NATS Streaming.
  * The `mcollectived` replacement that can run old ruby based agents
  * Everything using a Go implementation of the Choria Protocol so interop with the Ruby clients and servers are maintained.

This is heavily in progress and not really usable yet for the general public.

# Configuration

Sample configs are shown, subject to change

Running `choria broker run --config /path/to/broker.cfg` will start the various broker components, you can safely run the Middleware Broker, Federation Broker and Protocol Adapter all in the same process.

## NATS based Middleware Broker

This sets up a managed NATS instance, it's functionally equivalent to just running NATS standalone but it's easier to get going and with fewer settings to consider.

SSL is setup to be compatible with Choria - ie. uses the Puppet certificates etc.

```ini
# enables the middleware broker
plugin.choria.broker_network = true 

# port it listens on for clients, this is the default when not set
plugin.choria.network.client_port = 4222 

# port it listens on for other choria brokers for clustering purposes, this is the default when not set
plugin.choria.network.peer_port = 5222

# other brokers in a cluster
plugin.choria.network.peers = nats://choria1:5222, nats://choria2:5222, nats://choria3:5222

# username and password for cluster connections, no need for this typically, it would use CA validated TLS
# allowing only certs signed by your CA to connect
plugin.choria.network.peer_user = choria_cluster
plugin.choria.network.peer_password = s£cret

# enables the typical NATS stats/status port, default, set to 0 to disable
plugin.choria.network.monitor_port = 8222
```

## Federation Broker

The Federation Broker for now is configured using the exact same method as the Ruby one, it should be a drop in replacement.

```ini
# enables the federation broker
plugin.choria.broker_federation = true

# these settings are identical as the ruby one so I wont show them all
plugin.choria.broker_federation_cluster = production
plugin.choria.federation.instance = 1
```

## Protocol Adapter

The Protocol Adapter is a new feature that exist to adapt Choria traffic into other systems.  The initial use case is to receive all Registration data within a specific Collective and publish those into NATS Streaming.

I imagine a number of other scenarios:

  * Publishing registration data to Kafka, Lambda, Search systems, other CMDB
  * Updating PuppetDB facts more frequently than node runs to optimize discovery against it
  * Setting up generic listeners like `choria.adapter.elk` that can be used to receive replies and publish them into ELK.  You can do `mco rpc ... --reply-to choria.adapter.elk --nr` which would then not show the results to the user but instead publish them to ELK

Here we configure the NATS Streaming Adapter.  It listens for `request` messages from the old school MCollective Registration system and republish those messages into NATS Streaming where you can process them at a more leisurely pace and configure retention to your own needs.

```ini
# sets up a named adapter, you can run many of the same type
plugin.choria.adapters = discovery
plugin.choria.adapters.discovery.type = nats_stream

# configure this discovery adapter
# in this case the adapter does NATS->NATS Streaming so you need to configure both sides
# here is NATS Streaming
plugin.choria.adapter.discovery.stream.servers = stan1:4222,stan2:4222
plugin.choria.adapter.discovery.stream.clusterid = prod
plugin.choria.adapter.discovery.stream.topic = discovery # defaults to same as adapter name
plugin.choria.adapter.discovery.stream.workers = 10 # default

# here is the Collective side
plugin.choria.adapter.discovery.ingest.topic = mcollective.broadcast.agent.discovery
plugin.choria.adapter.discovery.ingest.protocol = request # or reply
plugin.choria.adapter.discovery.ingest.workers = 10 # default
```

## Choria Server

This will eventually replace `mcollectived`, for now all it can do is publish registration data.

You run it with `choria server run --config server.cfg`

Apart from all the usual stuff about identity, logfile etc, you can enable the new registration publisher like this, it just publishes registration data found in the file every 10 seconds:

```ini
registration = file_content
registerinterval = 10
registration_splay = true
plugin.choria.registration.file_content.data = /tmp/json_registration.json
```
