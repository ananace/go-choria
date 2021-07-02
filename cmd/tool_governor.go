package cmd

import "sync"

type tGovCommand struct {
	command
}

func (g *tGovCommand) Setup() (err error) {
	if tool, ok := cmdWithFullCommand("tool"); ok {
		g.cmd = tool.Cmd().Command("governor", "Distributed concurrency control system for Choria Streams").Alias("gov")
	}

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
