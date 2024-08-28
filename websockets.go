package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	e := createEngine(conn)

	if err != nil {
		log.Fatal("Serve http: ", err)
		return
	}

	stopEngine := make(chan struct{})
	go e.run(stopEngine)

    defer func() {
		// stop engine goroutine
		close(stopEngine)
		fmt.Println("stopping engine")
        conn.Close()
    }()
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("ReadMessage:", err)
            break
        }
		e.queueCommand(message)
    }
}