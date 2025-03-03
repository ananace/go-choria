// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/choria-io/go-choria/internal/util"
	agents "github.com/choria-io/go-choria/providers/agent/mcorpc/ddl/agent"
	"github.com/choria-io/go-choria/providers/data"
	dps "github.com/choria-io/go-choria/providers/data/ddl"
)

type pDocCommand struct {
	name     string
	jsonOnly bool
	markdown bool

	command
}

func (d *pDocCommand) Setup() (err error) {
	if tool, ok := cmdWithFullCommand("plugin"); ok {
		d.cmd = tool.Cmd().Command("doc", "Inspect plugin documentation")
		d.cmd.Arg("name", "Plugin to inspect").StringVar(&d.name)
		d.cmd.Flag("json", "Produce JSON output only").Short('j').BoolVar(&d.jsonOnly)
		d.cmd.Flag("markdown", "Produce Markdown output").Short('m').BoolVar(&d.markdown)
	}

	return nil
}

func (d *pDocCommand) Configure() error {
	return commonConfigure()
}

func (d *pDocCommand) Run(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	if d.name == "" {
		return d.showList()
	}

	return d.showPlugin()
}

func (d *pDocCommand) findAgent(name string) (*agents.DDL, error) {
	resolvers, err := c.DDLResolvers()
	if err != nil {
		return nil, err
	}

	log := c.Logger("ddl")

	for _, resolver := range resolvers {
		log.Infof("Resolving DDL agent/%s via %s", name, resolver)
		data, err := resolver.DDLBytes(ctx, "agent", name, c)
		if err == nil {
			return agents.NewFromBytes(data)
		}
	}

	return nil, fmt.Errorf("agent/%s ddl not found", name)
}

func (d *pDocCommand) agents() (map[string]*agents.DDL, error) {
	resolvers, err := c.DDLResolvers()
	if err != nil {
		return nil, err
	}

	found := map[string]struct{}{}
	res := make(map[string]*agents.DDL, len(found))
	log := c.Logger("ddl")

	for _, resolver := range resolvers {
		log.Infof("Resolving DDL names via %s", resolver)
		names, err := resolver.DDLNames(ctx, "agent", c)
		if err != nil {
			log.Warnf("Could not resolve DDL names using resolver %s: %v", resolver.String(), err)
			continue
		}

		for _, name := range names {
			found[name] = struct{}{}
		}
	}

	for name := range found {
		ddl, err := d.findAgent(name)
		if err != nil {
			return nil, err
		}
		res[ddl.Metadata.Name] = ddl
	}

	return res, nil
}

func (d *pDocCommand) data() (map[string]*dps.DDL, error) {
	dp, err := data.NewManager(context.Background(), c)
	if err != nil {
		return nil, err
	}

	found := dp.DDLs()
	res := make(map[string]*dps.DDL, len(found))
	for _, d := range found {
		res[d.Metadata.Name] = d
	}

	return res, nil
}

func (d *pDocCommand) showList() error {
	plugins := map[string]map[string]string{
		"agent": {},
		"data":  {},
	}

	var err error

	if !d.jsonOnly && cfg.Choria.RegistryClientCache != "" {
		fmt.Println("Accessing Choria Registry for service definitions")
		fmt.Println()
	}

	addls, _ := d.agents()
	for _, addl := range addls {
		plugins["agent"][addl.Metadata.Name] = addl.Metadata.Description
	}

	dddls, err := d.data()
	if err != nil {
		return err
	}
	for _, dddl := range dddls {
		plugins["data"][dddl.Metadata.Name] = dddl.Metadata.Description
	}

	if d.jsonOnly {
		out, err := json.MarshalIndent(plugins, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}

	fmt.Println("Known Plugins:")
	fmt.Println()
	fmt.Println("Agents:")
	fmt.Println()
	util.IterateStringsMap(plugins["agent"], func(k, v string) {
		fmt.Printf("%20s: %s\n", k, v)
	})

	fmt.Println()
	fmt.Println("Data Providers:")
	fmt.Println()
	util.IterateStringsMap(plugins["data"], func(k, v string) {
		fmt.Printf("%20s: %s\n", k, v)
	})
	fmt.Println()

	return nil
}

func (d *pDocCommand) showPlugin() error {
	ptype := "agent"
	pname := d.name
	if strings.Contains(d.name, "/") {
		parts := strings.Split(strings.ToLower(d.name), "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid plugn name %s", d.name)
		}
		ptype = parts[0]
		pname = parts[1]
	}

	switch ptype {
	case "agent":
		return d.renderAgent(pname)
	case "data":
		return d.renderData(pname)
	default:
		return fmt.Errorf("invalid plugin type %s", ptype)
	}
}

func (d *pDocCommand) renderAgent(agent string) error {
	ddl, err := d.findAgent(agent)
	if err != nil {
		return err
	}

	switch {
	case d.jsonOnly:
		return util.DumpJSONIndent(ddl)

	case d.markdown:
		out, err := ddl.RenderMarkdown()
		if err != nil {
			return err
		}
		fmt.Print(string(out))

	default:
		out, err := ddl.RenderConsole()
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}

	return nil
}

func (d *pDocCommand) renderData(data string) error {
	providers, err := d.data()
	if err != nil {
		return err
	}

	ddl, ok := providers[data]
	if !ok {
		return fmt.Errorf("unknown data provider %s", data)
	}

	switch {
	case d.jsonOnly:
		return util.DumpJSONIndent(ddl)

	case d.markdown:
		out, err := ddl.RenderMarkdown()
		if err != nil {
			return err
		}
		fmt.Print(string(out))

	default:
		out, err := ddl.RenderConsole()
		if err != nil {
			return err
		}
		fmt.Print(string(out))

	}

	return nil
}

func init() {
	cli.commands = append(cli.commands, &pDocCommand{})
}
