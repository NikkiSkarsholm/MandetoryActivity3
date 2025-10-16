package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var incommingMessages = make(chan proto.Message)
var users int = 0
var userIndexCounter int32 = 0
var userChannels = make(map[int32]chan proto.Message)

type ChitChatServer struct {
	proto.UnimplementedChitChatServer
	savedMessages []proto.Message
}

func main() {
	server := &ChitChatServer{}

	server.startServer()
}

func (s *ChitChatServer) startServer() {
	fmt.Println("Starting server")
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Server did not work")
	}

	proto.RegisterChitChatServer(grpcServer, s)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Server did not work")
	}
}

func (s *ChitChatServer) GetStream(clientID *proto.IdMessage, stream proto.ChitChat_GetStreamServer) (err error) {
	// Based on the Message object's lamport timestamp, determine which message to stream out first

	for {
		message := <-userChannels[clientID.GetId()]

		err := stream.Send(&message)
		if err != nil {
			return err
		}
	}
}

func (s *ChitChatServer) MessageToServer(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {

	for _, channel := range userChannels {
		channel <- *msg
	}

	return &proto.Empty{}, nil
}

func (s *ChitChatServer) UserJoins(ctx context.Context, joinMessage *proto.JoinMessage) (*proto.IdMessage, error) {
	fmt.Println("Participant '" + joinMessage.GetUsername() + "' joined chit chat at logical time")
	users++
	fmt.Println("There are currectly ", users, " manny users")

	id := userIndexCounter
	userIndexCounter++

	// give the user a channel
	userChannels[id] = make(chan proto.Message, 1)

	// send message to clients that the user joined
	test := "Participant '" + joinMessage.GetUsername() + "' joined chit chat at logical time"
	author := "server"
	Msg := proto.Message{Msg: test, Author: author}
	s.MessageToServer(context.Background(), &Msg)

	return &proto.IdMessage{Id: id}, nil
}

func (s *ChitChatServer) UserLeaves(ctx context.Context, leaveMessage *proto.LeaveMessage) (*proto.Empty, error) {
	fmt.Println("Participant '" + leaveMessage.GetUsername() + "' left chit chat at logical time")
	users--
	fmt.Println("There are currectly ", users, " manny users")

	fmt.Println("Deleting channel nr ", leaveMessage.GetId(), " number of channels are: ", len(userChannels))
	delete(userChannels, leaveMessage.GetId())
	fmt.Println("After deleting channel nr ", leaveMessage.GetId(), " number of channels are: ", len(userChannels))

	// send message to clients that the user left
	test := "Participant '" + leaveMessage.GetUsername() + "' left chit chat at logical time"
	author := "server"
	Msg := proto.Message{Msg: test, Author: author}
	s.MessageToServer(context.Background(), &Msg)

	return &proto.Empty{}, nil
}
