package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var engines = make(map[string]*engine)

type room struct {
	clients map[*client]bool

	join chan *client

	leave chan *client

	forward chan []byte
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func onDisconnect(e *engine) {
    fmt.Println("OnDisconnect called")
	e.running = false
	// e.socket.Close()
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("Serve http: ", err)
		return
	}

	params := req.URL.Query()
	roomName := params.Get("room_name")

	var engine *engine
	engine, exists := engines[roomName]
	if exists {
		fmt.Println("Engine already exists for room ", roomName)
		engine.running = true
		engine.socket = socket
		go engine.run()
	} else {
		fmt.Println("Creating new engine for room ", roomName)
		engine = startEngine(socket)
		engines[roomName] = engine
	}
	
	closeHandler := socket.CloseHandler()
	socket.SetCloseHandler(func(code int, text string) error {
		// Add your code here ...
		fmt.Println("WebSocket disconnected, stopping engine ", roomName)
		onDisconnect(engine)
		err := closeHandler(code, text)
		// ... or here.
		return err
	})

	// client := &client{
	// 	socket: socket,
	// 	receive: make(chan []byte, messageBufferSize),
	// 	room: r,
	// }
	// r.join <- client
	// defer func() { r.leave <- client }()
	// go client.write()
	// client.read()
}