package main

import (
	"log"
	"net/http"
)

func main () {
	r := newRoom()

	http.Handle("/ws", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}










// written by chapgpt but I think it works
// package main

// import (
//     "log"
//     "net/http"
//     "github.com/gorilla/websocket"
// )

// // Client represents a single WebSocket connection
// type Client struct {
//     hub *Hub
//     conn *websocket.Conn
//     send chan []byte
// }

// // Hub maintains the set of active clients and broadcasts messages to the clients
// type Hub struct {
//     register chan *Client
//     unregister chan *Client
//     clients map[*Client]bool
// }

// var upgrader = websocket.Upgrader{
//     ReadBufferSize:  1024,
//     WriteBufferSize: 1024,
// }

// func newHub() *Hub {
//     return &Hub{
//         register: make(chan *Client),
//         unregister: make(chan *Client),
//         clients: make(map[*Client]bool),
//     }
// }

// func (h *Hub) run() {
//     for {
//         select {
//         case client := <-h.register:
//             h.clients[client] = true
//             log.Println("Client registered")
//         case client := <-h.unregister:
//             if _, ok := h.clients[client]; ok {
//                 delete(h.clients, client)
//                 close(client.send)
//                 log.Println("Client unregistered")
//             }
//         }
//     }
// }

// func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
//     conn, err := upgrader.Upgrade(w, r, nil)
//     if err != nil {
//         log.Println(err)
//         return
//     }
//     client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
//     client.hub.register <- client // Register the client with the hub
// }

// func main() {
//     hub := newHub()
//     go hub.run()
//     http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//         serveWs(hub, w, r)
//     })
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }
