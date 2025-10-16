package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ChitChatClient
var userName string
var id int32

func main() {
	// make profile
	fmt.Println("Please enter your usename:")
	userName = readTerminal()

	// connecyting to server
	fmt.Println("Connectin to server ...")
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	client = proto.NewChitChatClient(conn)
	joinServer()
	fmt.Println("Connection established")

	// Print messeges using the message stream
	go messageStream(client)

	// read user inputs and send to server
	for {
		userInput := readTerminal()

		if userInput == "Quit" || userInput == "quit" || userInput == "Q" || userInput == "q" {
			leaveServer()
			return
		}

		messageServer(userInput)
	}
}

func messageServer(message string) {
	_, err := client.MessageToServer(context.Background(), &proto.Message{Msg: message, Author: userName})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
}

func messageStream(client proto.ChitChatClient) {
	stream, err := client.GetStream(context.Background(), &proto.IdMessage{Id: id})
	if err != nil {
		log.Fatalf("could not get stream: %v", err)
	}

	for {
		msg, _ := stream.Recv()
		printMessage(msg)
	}
}

func printMessage(msg *proto.Message) {
	fmt.Println(msg.Author + ": " + msg.Msg)
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)

	//fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF
	text = strings.Replace(text, "\r", "", -1) // convert CRLF to LF

	return text
}
func joinServer() {
	IdMessage, err := client.UserJoins(context.Background(), &proto.JoinMessage{Username: userName})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	id = IdMessage.GetId()
	fmt.Println("You have joined the server as ", userName, " with the id: ", id)
}

func leaveServer() {
	_, err := client.UserLeaves(context.Background(), &proto.LeaveMessage{Username: userName, Id: id})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
}
