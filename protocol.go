package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type message struct {
	Type      string      `json:"t"`
	TimeStamp time.Time   `json:"d"`
	Payload   interface{} `json:"p"`
}

func (m *message) UnmarshalJSON(data []byte) (err error) {
	fmt.Printf("unmashaling json \n%s\n", data)
	raw := &struct {
		T string          `json:"t"`
		D time.Time       `json:"d"`
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
	j, _ := json.Marshal(m.Payload)
	fmt.Printf("payload is json \n%s\n", j)
	return json.Unmarshal(raw.P, m.Payload)
}

type patchPayload struct {
	Author  string                 `json:"a"`
	Patches []diffmatchpatch.Patch `json:"p"`
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
