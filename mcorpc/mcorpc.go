// Package mcorpc provides a compatibility layer between Choria and
// legacy MCollective SimpleRPC Agents
//
// Agents can be written in the Go language, compiled into the binaries
// and be interacted with from the ruby MCollective client.
//
// It's planned to provide a backward compatible interface so that old
// ruby agents, authorization and auditing will be usable inside the
// Choria daemon via a shell-out mechanism
package mcorpc

import (
	"encoding/json"
	"fmt"
)

// StatusCode is a reply status as defined by MCollective SimpleRPC - integers 0 to 5
//
// See the constants OK, RPCAborted, UnknownRPCAction, MissingRPCData, InvalidRPCData and UnknownRPCError
type StatusCode uint8

const (
	// OK is the reply status when all worked
	OK = StatusCode(iota)

	// Aborted is status for when the action could not run, most failures in an action should set this
	Aborted

	// UnknownAction is the status for unknown actions requested
	UnknownAction

	// MissingData is the status for missing input data
	MissingData

	// InvalidData is the status for invalid input data
	InvalidData

	// UnknownError is the status general failures in agents should set when things go bad
	UnknownError
)

// Reply is the reply data as stipulated by MCollective RPC system.  The Data
// has to be something that can be turned into JSON using the normal Marshal system
type Reply struct {
	Statuscode StatusCode  `json:"statuscode"`
	Statusmsg  string      `json:"statusmsg"`
	Data       interface{} `json:"data"`
}

// Request is a request as defined by the MCollective RPC system
// NOTE: input arguments not yet handled
type Request struct {
	Agent  string          `json:"agent"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// ParseRequestData parses the request parameters received from the client into a target structure
//
// Example used in a action:
//
//   var rparams struct {
//      Package string `json:"package"`
//   }
//
//   if !mcorpc.ParseRequestData(&rparams, req, reply) {
//     // the function already set appropriate errors on reply
//	   return
//   }
//
//   // do stuff with rparams.Package
func ParseRequestData(target interface{}, request *Request, reply *Reply) bool {
	err := json.Unmarshal(request.Data, target)
	if err != nil {
		reply.Statuscode = InvalidData
		reply.Statusmsg = fmt.Sprintf("Could not parse request data for %s#%s: %s", request.Agent, request.Action, err)

		return false
	}

	return true
}
