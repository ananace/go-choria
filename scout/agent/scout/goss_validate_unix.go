// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aelsabbahy/goss"
	gossoutputs "github.com/aelsabbahy/goss/outputs"
	gossutil "github.com/aelsabbahy/goss/util"
	"github.com/choria-io/go-choria/inter"

	"github.com/choria-io/go-choria/providers/agent/mcorpc"
)

type GossValidateRequest struct {
	File string `json:"file"`
	Vars string `json:"vars"`
}

type GossValidateResponse struct {
	Failures int                                `json:"failures"`
	Results  []gossoutputs.StructuredTestResult `json:"results"`
	Runtime  float64                            `json:"runtime"`
	Success  int                                `json:"success"`
	Summary  string                             `json:"summary"`
	Tests    int                                `json:"tests"`
}

func gossValidateAction(_ context.Context, req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, _ inter.ConnectorInfo) {
	resp := &GossValidateResponse{Results: []gossoutputs.StructuredTestResult{}}
	reply.Data = resp

	args := &GossValidateRequest{}
	if !mcorpc.ParseRequestData(args, req, reply) {
		return
	}

	var out bytes.Buffer

	opts := []gossutil.ConfigOption{
		gossutil.WithMaxConcurrency(1),
		gossutil.WithResultWriter(&out),
		gossutil.WithSpecFile(args.File),
	}

	if args.Vars != "" {
		opts = append(opts, gossutil.WithVarsFile(args.File))
	}

	cfg, err := gossutil.NewConfig(opts...)
	if err != nil {
		abort(fmt.Sprintf("Could not create Goss config: %s", err), reply)
		return
	}

	_, err = goss.Validate(cfg, time.Now())
	if err != nil {
		abort(fmt.Sprintf("Could not validate: %s", err), reply)
		return
	}

	res := &gossoutputs.StructuredOutput{}
	err = json.Unmarshal(out.Bytes(), res)
	if err != nil {
		abort(fmt.Sprintf("Could not parse goss results: %s", err), reply)
		return
	}

	resp.Results = res.Results
	resp.Summary = res.SummaryLine
	resp.Failures = res.Summary.Failed
	resp.Runtime = res.Summary.TotalDuration.Seconds()
	resp.Success = res.Summary.TestCount - res.Summary.Failed
	resp.Tests = res.Summary.TestCount
}
