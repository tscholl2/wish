package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type room struct {
	// lock for map access
	sync.RWMutex

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
	r.RLock()
	for c := range r.connections {
		if err := c.send(websocket.TextMessage, msg); err != nil {
			c.close()
		}
	}
	r.RUnlock()
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
