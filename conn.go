package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newConnection(w http.ResponseWriter, r *http.Request) (*connection, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := &connection{
		ws: ws,
	}
	return c, nil
}

// warning: make sure only one goroutine calls this at a time
func (c *connection) read() (int, []byte, error) {
	return c.ws.ReadMessage()
}

// warning: make sure only one goroutine calls this at a time
// also if there is an error, the source says it immediately returns
// that error on subsequent calls, so it's safe
func (c *connection) send(msgType int, msg []byte) error {
	return c.ws.WriteMessage(msgType, msg)
}

// anyone can call this whenever
func (c *connection) close() error {
	c.ws.WriteMessage(websocket.CloseMessage, []byte{})
	return c.ws.Close()
}
