// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package agents

import (
	"github.com/choria-io/go-choria/build"
	"github.com/choria-io/go-choria/inter"
	"github.com/choria-io/go-choria/plugin"
)

// ChoriaPlugin produces the plugin for choria
func ChoriaPlugin() plugin.Pluggable {
	return &Provider{}
}

// PluginInstance implements plugin.Pluggable
func (p *Provider) PluginInstance() interface{} {
	return p
}

// PluginVersion implements plugin.Pluggable
func (p *Provider) PluginVersion() string {
	return build.Version
}

// PluginName implements plugin.Pluggable
func (p *Provider) PluginName() string {
	return "Golang MCollective Agent Compatibility"
}

// PluginType implements inter.Pluggable
func (p *Provider) PluginType() inter.PluginType {
	return inter.AgentProviderPlugin
}
