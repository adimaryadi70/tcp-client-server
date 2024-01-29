// client.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:9898", "WebSocket server address")

func main() {
	flag.Parse()
	log.Printf("Connecting to %s", *addr)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+*addr+"/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	done := make(chan struct{})

	// Start a goroutine to receive messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Received message from server: %s\n", message)
		}
	}()

	// Start a goroutine to read user input and send it to the server
	go func() {
		defer close(done)
		for {
			var message string
			fmt.Print("Enter message: ")
			fmt.Scanln(&message)

			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Wait for the interrupt signal or done signal
	select {
	case <-interrupt:
	case <-done:
	}
}
