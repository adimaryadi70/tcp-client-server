package main

import (
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/castaneai/grpc-broadcast-example"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type User struct {
	chat pb.ChatRoom_ChatServer
	move pb.ChatRoom_ChatServer
}

type server struct {
	clients map[string]pb.ChatRoom_ChatServer
	mu      sync.RWMutex
}

func (s *server) addClient(uid string, srv pb.ChatRoom_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[uid] = srv
}

func (s *server) removeClient(uid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, uid)
}

func (s *server) getClients() []pb.ChatRoom_ChatServer {
	var cs []pb.ChatRoom_ChatServer

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.clients {
		cs = append(cs, c)
	}
	return cs
}

func (s *server) Chat(srv pb.ChatRoom_ChatServer) error {
	uid := uuid.Must(uuid.NewRandom()).String()
	log.Println("id User:", uid)

	s.addClient(uid, srv)

	defer s.removeClient(uid)

	defer func() {
		if err := recover(); err != nil {
			log.Println("error :", err)
			os.Exit(1)
		}
	}()

	for {
		response, err := srv.Recv()
		if err != nil {
			log.Printf("recv err: %v", err)
			break
		}
		log.Printf("broadcast: %s", response.Message)
		for _, data := range s.getClients() {
			if err := data.Send(&pb.ChatResponse{Message: response.Message}); err != nil {
				log.Printf("broadcast err: %v", err)
			}
		}
	}
	return nil
}

func main() {
	address := ":9898"
	listen, err := net.Listen("tcp", address)
	log.Println("Server Listen:", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatRoomServer(s, &server{
		clients: make(map[string]pb.ChatRoom_ChatServer),
		mu:      sync.RWMutex{},
	})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
