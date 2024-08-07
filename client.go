package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn

	// receive is a  channel which receives messages from other clients
	receive chan []byte

	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	fmt.Println("reading...")
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}