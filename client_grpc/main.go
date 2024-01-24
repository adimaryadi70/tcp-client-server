package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/adimaryadi70/proto/example"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9898", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error to Connect:", err)
	}
	defer conn.Close()

	client := pb.NewMyServiceClient(conn)

	// req := &pb.MyRequest{
	// 	Data: "Tester Kirim GRPC",
	// }

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Message ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		req := &pb.MyRequest{
			Data: input,
		}

		ctx := context.Background()
		res, err := client.MyMethod(ctx, req)
		if err != nil {
			log.Fatal("Failed to call ", err)
		}
		fmt.Printf("Response from server: %s\n", res.Result)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}
