// Copyright (c) 2019-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/providers/agent/mcorpc/ddl/common"

	addl "github.com/choria-io/go-choria/providers/agent/mcorpc/ddl/agent"
)

func (g *Generator) templFSnakeToCamel(v string) string {
	parts := strings.Split(v, "_")
	out := []string{}
	for _, s := range parts {
		out = append(out, strings.Title(s))
	}

	return strings.Join(out, "")
}

func (g *Generator) templFSnakeToCamelUnexported(v string) string {
	parts := strings.Split(v, "_")
	out := []string{}
	for i, s := range parts {
		if i == 0 {
			out = append(out, strings.ToLower(s))
		} else {
			out = append(out, strings.Title(s))
		}
	}

	return strings.Join(out, "")
}

func (g *Generator) templFChoriaTypeToValOfType(v string) string {
	switch v {
	case "string", "list":
		return "val.(string)"
	case "integer":
		return "val.(int64)"
	case "number", "float":
		return "val.(float64)"
	case "boolean":
		return "val.(bool)"
	case "hash":
		return "val.(map[string]interface{})"
	case "array":
		return "val.([]interface{})"
	default:
		return "val.(interface{})"
	}
}

func (g *Generator) templFChoriaRequiredInputsToFuncArgs(act *addl.Action) string {
	inputs := g.optionalInputSelect(act, false)
	var parts []string

	names := []string{}
	for name := range inputs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		input := inputs[name]
		goType := g.templFChoriaTypeToGo(input.Type)
		parts = append(parts, fmt.Sprintf("input%s %s", g.templFSnakeToCamel(name), goType))
	}

	return strings.Join(parts, ", ")
}

func (g *Generator) templFChoriaTypeToGo(v string) string {
	switch v {
	case "string", "list":
		return "string"
	case "integer":
		return "int64"
	case "number", "float":
		return "float64"
	case "boolean":
		return "bool"
	case "hash":
		return "map[string]interface{}"
	case "array":
		return "[]interface{}"
	default:
		return "interface{}"
	}
}

func (g *Generator) templFChoriaOptionalInputs(act *addl.Action) map[string]*common.InputItem {
	return g.optionalInputSelect(act, true)
}

func (g *Generator) templFChoriaRequiredInputs(act *addl.Action) map[string]*common.InputItem {
	return g.optionalInputSelect(act, false)
}

func (g *Generator) templFBase64Encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func (g *Generator) templFGeneratedWarning() string {
	meta := g.agent.DDL.Metadata
	return fmt.Sprintf(`// generated code; DO NOT EDIT"
//
// Client for Choria RPC Agent '%s' Version %s generated using Choria version %s`, meta.Name, meta.Version, choria.BuildInfo().Version())
}

func (g *Generator) funcMap() template.FuncMap {
	return template.FuncMap{
		"GeneratedWarning":               g.templFGeneratedWarning,
		"Base64Encode":                   g.templFBase64Encode,
		"Capitalize":                     strings.Title,
		"ToLower":                        strings.ToLower,
		"SnakeToCamel":                   g.templFSnakeToCamel,
		"SnakeToCamelUnexported":         g.templFSnakeToCamelUnexported,
		"ChoriaRequiredInputs":           g.templFChoriaRequiredInputs,
		"ChoriaOptionalInputs":           g.templFChoriaOptionalInputs,
		"ChoriaRequiredInputsToFuncArgs": g.templFChoriaRequiredInputsToFuncArgs,
		"ChoriaTypeToGoType":             g.templFChoriaTypeToGo,
		"ChoriaTypeToValOfType":          g.templFChoriaTypeToValOfType,
	}
}
