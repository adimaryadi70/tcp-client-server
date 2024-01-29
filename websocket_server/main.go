// server.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	clients    = make(map[*client]bool)
	register   = make(chan *client)
	unregister = make(chan *client)
	broadcast  = make(chan []byte)
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{conn: conn, send: make(chan []byte, 256)}
	register <- client

	defer func() {
		unregister <- client
		conn.Close()
	}()

	for {
		_, p, err := conn.ReadMessage()
		log.Print("message", p)
		if err != nil {
			log.Println(err)
			return
		}
		broadcast <- p
	}
}

func handleMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			delete(clients, client)
			close(client.send)
		case message := <-broadcast:
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func handleClient(client *client) {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	port := ":9898"
	go handleMessages()
	log.Print("Listener Port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
