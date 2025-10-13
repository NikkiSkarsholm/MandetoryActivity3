package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := proto.NewChitChatClient(conn)

	_, err = client.MessageToServer(context.Background(), &proto.Message{Msg: "Hello World", Author: "Karam"})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	//go messageStream(client) // Connect to server stream

	//readTerminal()

	go readTerminal2()

	time.Sleep(10 * time.Second)
}

func messageStream(client proto.ChitChatClient) {
	stream, err := client.GetStream(context.Background(), &proto.Empty{})
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

/*
func createMessage() proto.Message {
}
*/

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input message: ")
	fmt.Println("---------------------")

	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF

	return text
}

func readTerminal2() {
	var input []byte
	reader := bufio.NewReader(os.Stdin)
	// disables input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// append each character that gets typed to input slice
	for {
		fmt.Println("-> " + string(input))
		b, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		input = append(input, b)
	}
}
