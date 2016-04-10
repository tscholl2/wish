package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type message struct {
	Type      string      `json:"t"`
	TimeStamp time.Time   `json:"d"`
	Payload   interface{} `json:"p"`
}

func (m *message) UnmarshalJSON(data []byte) (err error) {
	raw := &struct {
		T string          `json:"t"`
		D time.Time       `json:"d"`
		P json.RawMessage `json:"p"`
	}{}
	if err = json.Unmarshal(data, raw); err != nil {
		return
	}
	m.Type = raw.T
	m.TimeStamp = raw.D
	switch raw.T {
	case "p":
		m.Payload = new(patchPayload)
	case "s":
		m.Payload = new(snapshotPayload)
	default:
		return fmt.Errorf("unknown message type: %+v", m)
	}
	return json.Unmarshal(raw.P, m.Payload)
}

type patchPayload struct {
	Author  string  `json:"a"`
	Patches []patch `json:"p"`
}

type patch struct {
	Start int    `json:"1"`
	End   int    `json:"2"`
	Text  string `json:"s"`
}

type snapshotPayload struct {
	Text string `json:"t"`
}

func newSnapshotMessage(s string) *message {
	return &message{
		Type:      "s",
		TimeStamp: time.Now(),
		Payload: snapshotPayload{
			Text: s,
		},
	}
}
