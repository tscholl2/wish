package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type message struct {
	Type    string      `json:"t"`
	Payload interface{} `json:"p"`
}

func (m *message) UnmarshalJSON(data []byte) (err error) {
	raw := &struct {
		T string          `json:"t"`
		P json.RawMessage `json:"p"`
	}{}
	if err = json.Unmarshal(data, raw); err != nil {
		return
	}
	m.Type = raw.T
	switch raw.T {
	case "p":
		m.Payload = patchPayload{}
	case "s":
		m.Payload = snapshotPayload{}
	default:
		return fmt.Errorf("unknown message type: %+v", m)
	}
	return json.Unmarshal(raw.P, m.Payload)
}

type patchPayload struct {
	Client    string                 `json:"id"`
	TimeStamp time.Time              `json:"d"`
	Patches   []diffmatchpatch.Patch `json:"p"`
}

type snapshotPayload struct {
	TimeStamp time.Time `json:"d"`
	Text      string    `json:"t"`
}
