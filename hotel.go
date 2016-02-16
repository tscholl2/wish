package main

import "fmt"

type roomRequest struct {
	c    *connection
	name string
}

type hotel struct {
	// set of rooms in use
	rooms map[string]*room

	// Put a connection into  a room
	// and create one if it doesn't exist
	enter chan roomRequest

	// Remove a connection from a room
	// and deletes the room if no one is left in it
	leave chan roomRequest
}

func (H *hotel) run() {
	for {
		select {
		case req := <-H.enter:
			if _, ok := H.rooms[req.name]; !ok {
				H.rooms[req.name] = newRoom(req.name)
			}
			r := H.rooms[req.name]
			req.c.send(newSnapshotMessage(r.text.snapshot))
			r.Lock()
			r.connections[req.c] = true
			r.Unlock()
			go func(req roomRequest) {
				for {
					msg, err := req.c.read()
					fmt.Printf("got msg from %p, size ~ %d\n", req.c, len(fmt.Sprintf("%+v", msg)))
					if err != nil {
						H.leave <- req
						break
					} else {
						r.inbox <- msg
					}
				}
			}(req)
		case req := <-H.leave:
			r := H.rooms[req.name]
			r.Lock()
			delete(r.connections, req.c)
			r.Unlock()
			if len(r.connections) == 0 {
				delete(H.rooms, req.name)
				r.done <- struct{}{}
			}
			req.c.close()
		}
	}
}

func newHotel() *hotel {
	H := &hotel{
		rooms: make(map[string]*room),
		enter: make(chan roomRequest),
		leave: make(chan roomRequest),
	}
	go H.run()
	return H
}
