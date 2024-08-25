package main

import (
	"log"
	"net/http"
)

func main () {
	r := newRoom()

	http.Handle("/ws", r)


	go r.run()
	// startEngine()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}







// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"github.com/gorilla/websocket"
// )

// var upgrader2 = websocket.Upgrader{
// 	// CheckOrigin allows all connections by default.
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade initial GET request to a websocket
// 	ws, err := upgrader2.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer ws.Close()

// 	for {
// 		fmt.Println("waiting for message")
// 		// Read message from client
// 		_, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		// Print the received message
// 		fmt.Printf("Received: %s\n", msg)

// 		// Write message back to client
// 		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 	}
// 	fmt.Println("end of handle connections")
// }

// func main() {
// 	http.HandleFunc("/ws", handleConnections)

// 	fmt.Println("Starting server on :8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		fmt.Println(err)
// 	}
// }



// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// var upgrader2 = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func main() {
// 	http.HandleFunc("/ws", handleConnections)
// 	log.Println("HTTP server started on :8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader2.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer ws.Close()

// 	done := make(chan struct{})

// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				// WebSocket closed, but goroutine continues running
// 				fmt.Println("Goroutine continues running after WebSocket closes.")
// 				for {
// 					fmt.Println("Still running...")
// 					time.Sleep(2 * time.Second)
// 				}
// 			}
// 		}
// 	}()

// 	for {
// 		_, _, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println("WebSocket closed:", err)
// 			close(done) // Signal to the goroutine that the websocket is closed
// 			break
// 		}
// 	}
// }
