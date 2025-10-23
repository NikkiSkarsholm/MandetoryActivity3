package main

import (
	lamportclock "ChitChat/General"
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ChitChatClient
var clock lamportclock.SafeClock
var userName string
var id int32

func main() {
	// make profile
	clock.Iterate()
	log.Println("Please enter your username:")
	userName = readTerminal()

	// connecting to server
	log.Println("Connecting to server ...")
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	clock.Iterate()
	client = proto.NewChitChatClient(conn)
	joinServer()
	log.Println("Connection established")

	// Print messages using the message stream
	go messageStream(client)

	// read user inputs and send to server
	for {
		userInput := readTerminal()

		if len(userInput) > 128 || len(userInput) <= 0 || !utf8.ValidString(userInput) {
			log.Println("[WARNING] Could not send message. Message must not be empty, be valid in utf8 and have maximum length of 128 characters")
			continue
		}

		if userInput == "Quit" || userInput == "quit" || userInput == "Q" || userInput == "q" {
			leaveServer()
			return
		}

		messageServer(userInput)
	}
}

func messageServer(message string) {
	TimeMessage, err := client.MessageToServer(context.Background(), &proto.Message{Msg: message, Author: userName, LamportTimestamp: clock.GetTime(), Id: id})
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	clock.MatchTime(TimeMessage.GetLamportTimestamp())
}

func messageStream(client proto.ChitChatClient) {
	stream, err := client.GetStream(context.Background(), &proto.IdMessage{Id: id, LamportTimestamp: clock.GetTime()})
	if err != nil {
		log.Fatalf("Could not get stream: %v", err)
	}
	clock.Iterate()

	for {
		msg, _ := stream.Recv()
		clock.MatchTime(msg.GetLamportTimestamp())
		printMessage(msg)
	}
}

func printMessage(msg *proto.Message) {
	log.Println(msg.Author+": "+msg.Msg+". Client's current time: ", clock.GetTime())
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF
	text = strings.Replace(text, "\r", "", -1) // convert CRLF to LF
	clock.Iterate()
	return text
}

func joinServer() {
	IdMessage, err := client.UserJoins(context.Background(), &proto.JoinMessage{Username: userName, LamportTimestamp: clock.GetTime()})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	id = IdMessage.GetId()
	clock.MatchTime(IdMessage.GetLamportTimestamp())
	log.Println("You have joined the server as ", userName, " with the id: ", id)
}

func leaveServer() {
	TimeMessage, err := client.UserLeaves(context.Background(), &proto.LeaveMessage{Username: userName, Id: id, LamportTimestamp: clock.GetTime()})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	clock.MatchTime(TimeMessage.GetLamportTimestamp())
	log.Println("You have left the server at logical time ", clock.GetTime())
}
