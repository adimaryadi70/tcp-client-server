package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/castaneai/grpc-broadcast-example"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9898", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error to Connect:", err)
	}
	defer conn.Close()

	c := pb.NewChatRoomClient(conn)

	clientName := os.Args[0]

	ctx := context.Background()

	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			message := "test" + clientName
			if err := stream.SendMsg(&pb.ChatRequest{Message: message}); err != nil {
				log.Fatal(err)
			}

			log.Printf("Send: %s", message)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("recv: %s", resp.Message)
	}
}
