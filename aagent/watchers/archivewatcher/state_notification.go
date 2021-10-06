// Copyright (c) 2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"encoding/json"
	"fmt"

	"github.com/choria-io/go-choria/aagent/watchers/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// StateNotification describes the current state of the watcher
// described by io.choria.machine.watcher.exec.v1.state
type StateNotification struct {
	event.Event

	Source          string `json:"source"`
	Creates         string `json:"creates"`
	PreviousOutcome string `json:"previous_outcome"`
	PreviousRunTime int64  `json:"previous_run_time"`
}

// CloudEvent creates a CloudEvent from the state notification
func (s *StateNotification) CloudEvent() cloudevents.Event {
	return s.Event.CloudEvent(s)
}

// JSON creates a JSON representation of the notification
func (s *StateNotification) JSON() ([]byte, error) {
	return json.Marshal(s.CloudEvent())
}

// String is a string representation of the notification suitable for printing
func (s *StateNotification) String() string {
	return fmt.Sprintf("%s %s#%s source: %s, previous: %s ran: %.3fs", s.Identity, s.Machine, s.Name, s.Source, s.PreviousOutcome, float64(s.PreviousRunTime)/1000000000)
}
