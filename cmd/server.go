// Copyright (c) 2017-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"sync"
)

type serverCommand struct {
	command
}

type serverRunCommand struct {
	command

	serviceHost      bool
	disableTLS       bool
	disableTLSVerify bool
	pidFile          string
}

// server
func (b *serverCommand) Setup() (err error) {
	b.cmd = cli.app.Command("server", "Choria Server")
	b.cmd.Flag("config", "Config file to use").PlaceHolder("FILE").StringVar(&configFile)

	return
}

func (b *serverCommand) Run(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	return
}

func (b *serverCommand) Configure() error {
	cfg.DisableSecurityProviderVerify = true

	return nil
}

func init() {
	cli.commands = append(cli.commands, &serverCommand{})
}
