// Copyright (c) 2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import "sync"

type tGovCommand struct {
	command
}

func (g *tGovCommand) Setup() (err error) {
	g.cmd = cli.app.Command("governor", "Distributed concurrency control system for Choria Streams").Alias("gov")
	g.cmd.Flag("config", "Config file to use").PlaceHolder("FILE").StringVar(&configFile)

	return nil
}

func (g *tGovCommand) Configure() error {
	return nil
}

func (g *tGovCommand) Run(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	return nil
}

func init() {
	cli.commands = append(cli.commands, &tGovCommand{})
}
