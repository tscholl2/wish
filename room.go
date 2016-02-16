package main

import "sync"

type room struct {
	// lock for map access
	sync.Mutex

	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	inbox chan *message

	// text storage
	text *text

	// done signal
	done chan struct{}
}

func (r *room) run() {
	for {
		select {
		case msg := <-r.inbox:
			switch msg.Payload.(type) {
			case patchPayload:
				r.text.update(msg.Payload.(patchPayload).Patches)
				r.send(msg)
			}
		case <-r.done:
			break
		}
	}
}

func (r *room) send(msg *message) {
	r.Lock()
	for c := range r.connections {
		if err := c.send(msg); err != nil {
			c.close()
		}
	}
	r.Unlock()
}

func newRoom(name string) *room {
	r := &room{
		connections: make(map[*connection]bool),
		inbox:       make(chan *message),
		done:        make(chan struct{}),
	}
	r.text = newText()
	go r.run()
	return r
}
