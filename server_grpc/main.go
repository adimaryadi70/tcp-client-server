package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/adimaryadi70/proto/example"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) MyMethod(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	result := fmt.Sprintf("Received: %s", req.Data)
	log.Println("Message ", req.Data)
	return &pb.MyResponse{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9898")
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMyServiceServer(grpcServer, &server{})

	fmt.Println("gRPC server is running on :9898")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
