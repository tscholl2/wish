package main

import (
	"fmt"
	"time"
)

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
}

func (H *hotel) run() *hotel {
	if H.rooms != nil || H.enter != nil {
		panic("Hotel can only run once")
	}
	H.rooms = make(map[string]*room)
	H.enter = make(chan roomRequest)
	go func(H *hotel) {
		for {
			select {
			case req := <-H.enter:
				if _, ok := H.rooms[req.name]; !ok {
					H.rooms[req.name] = new(room).run()
					fmt.Printf("new room: %s\n", req.name)
				}
				r := H.rooms[req.name]
				r.add(req.c)
			case <-time.Tick(10 * time.Second):
				for name, r := range H.rooms {
					r.Lock()
					if len(r.connections) == 0 {
						delete(H.rooms, name)
						fmt.Printf("kill room: %s\n", name)
					}
					r.Unlock()
					r.done <- struct{}{}
				}
			}
		}
	}(H)
	return H
}
