package main

import (
	"encoding/json"
	"fmt"
	"time"
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
	Start   int       `json:"1"`
	End     int       `json:"2"`
	NewText string    `json:"n"`
	OldText string    `json:"o"`
	Time    time.Time `json:"t"`
	ID      string    `json:"id"`
}

type snapshotPayload struct {
	Timestamp string `json:"d"`
	Text      string `json:"t"`
}

func newSnapshotMessage(s string) *message {
	d, _ := time.Now().UTC().MarshalJSON()
	return &message{
		Type:    "s",
		Payload: snapshotPayload{string(d), s},
	}
}
