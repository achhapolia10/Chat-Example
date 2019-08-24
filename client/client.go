package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/achhapolia10/chatExample/chatpb"

	"google.golang.org/grpc"
)

func main() {
	w := make(chan struct{})

	clientConn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error in Connecting to server: %v", err)
	}
	client := chatpb.NewChatServiceClient(clientConn)

	uname := handleLogin(client)
	go sendMessages(client, uname)

	<-w
}

func handleLogin(c chatpb.ChatServiceClient) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter Username: ")
	uname, _ := reader.ReadString('\n')
	uname = uname[0 : len(uname)-1]
	l, _ := c.Login(context.Background(), &chatpb.LoginRequest{Username: uname})
	fmt.Println("User Logged in")
	reciever, _ := c.StartReciveing(context.Background(), l)
	go recieveMessages(reciever)
	return uname
}

func sendMessages(c chatpb.ChatServiceClient, uname string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Enter Message: ")
		message, _ := reader.ReadString('\n')
		message = message[0 : len(message)-1]

		fmt.Printf("Enter Recipient: ")
		r, _ := reader.ReadString('\n')
		r = r[0 : len(r)-1]
		c.SendMessage(context.Background(), &chatpb.Messages{
			Sender:   uname,
			Reciever: r,
			Message:  message,
		})
	}
}

func recieveMessages(r chatpb.ChatService_StartReciveingClient) {
	for {
		message, _ := r.Recv()
		fmt.Printf("%s : %s", message.GetSender(), message.GetMessage())
	}
}
