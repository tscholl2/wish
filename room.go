package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type room struct {
	// lock for map access
	sync.Mutex

	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	inbox chan []byte

	// done signal
	done chan struct{}
}

func (r *room) run() {
	for {
		select {
		case msg := <-r.inbox:
			r.send(websocket.TextMessage, msg)
		case <-r.done:
			break
		}
	}
}

func (r *room) send(msgType int, msg []byte) {
	r.Lock()
	for c := range r.connections {
		if err := c.send(msgType, msg); err != nil {
			c.close()
		}
	}
	r.Unlock()
}

func newRoom(name string) *room {
	r := &room{
		connections: make(map[*connection]bool),
		inbox:       make(chan []byte),
		done:        make(chan struct{}),
	}
	go r.run()
	return r
}
