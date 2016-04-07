package main

import (
	"fmt"
	"sync"
)

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
				fmt.Println("room got a msg!")
				switch msg.Payload.(type) {
				case *patchPayload:
					fmt.Println("its a patch!")
					r.text.update(msg.Payload.(*patchPayload).Patches)
					fmt.Printf("new text = \n%s\n", r.text.snapshot)
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
	fmt.Printf("new connection! %p\n", c)
	r.Lock()
	r.connections[c] = struct{}{}
	c.send(newSnapshotMessage(r.text.snapshot))
	go func(c *connection, r *room) {
		for {
			msg, err := c.read()
			fmt.Printf("new msg from %p\n%s\n", c, msg)
			if err != nil {
				fmt.Printf("err msg from %p, removing\n", c)
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
	fmt.Printf("remove connection! %p\n", c)
	r.Lock()
	delete(r.connections, c)
	r.Unlock()
}
