package main

import "sync"

type room struct {
	// lock for map access
	sync.Mutex

	// Registered connections.
	connections map[*connection]struct{}

	// Inbound messages from the connections.
	inbox chan *message

	// text storage
	text *text

	// done signal
	done chan struct{}
}

func (r *room) run() *room {
	if r.connections != nil || r.inbox != nil || r.done != nil || r.text != nil {
		panic("rooms can only run once")
	}
	r.connections = make(map[*connection]struct{})
	r.inbox = make(chan *message)
	r.done = make(chan struct{})
	r.text = new(text)
	go func(r *room) {
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
	}(r)
	return r
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

func (r *room) add(c *connection) {
	r.Lock()
	r.connections[c] = struct{}{}
	r.send(newSnapshotMessage(r.text.snapshot))
	go func(c *connection, r *room) {
		for {
			msg, err := c.read()
			if err != nil {
				r.remove(c)
				break
			} else {
				r.inbox <- msg
			}
		}
	}(c, r)
	r.Unlock()
}

func (r *room) remove(c *connection) {
	r.Lock()
	delete(r.connections, c)
	r.Unlock()
}
