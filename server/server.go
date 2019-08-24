package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/achhapolia10/chatExample/chatpb"
)

type server struct {
	uRecChan map[string]chan *chatpb.Messages
}

func main() {
	log.Printf("Starting Chat Server on PORT 8080")

	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Error in Listening to port 8080 : %v", err)
	}

	gs := grpc.NewServer()

	s := &server{
		uRecChan: make(map[string]chan *chatpb.Messages),
	}

	chatpb.RegisterChatServiceServer(gs, s)

	if err := gs.Serve(l); err != nil {
		log.Fatalf("Error in serving the server: %v", err)
	}

}

func (s *server) Login(c context.Context, req *chatpb.LoginRequest) (*chatpb.LoginResponse, error) {
	return &chatpb.LoginResponse{
		LoginResult: req.GetUsername(),
	}, nil
}

func (s *server) StartReciveing(req *chatpb.LoginResponse, res chatpb.ChatService_StartReciveingServer) error {
	s.uRecChan[req.GetLoginResult()] = make(chan *chatpb.Messages)
	for {

		m := <-s.uRecChan[req.GetLoginResult()]
		res.Send(m)
	}
	return nil
}
func (s *server) SendMessage(c context.Context, req *chatpb.Messages) (*chatpb.SendMessageResponse, error) {
	s.uRecChan[req.Reciever] <- req
	return &chatpb.SendMessageResponse{
		Status: "send",
	}, nil
}
