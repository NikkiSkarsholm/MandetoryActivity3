package main

import (
	proto "ChitChat/grpc"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChitChatServer struct {
	proto.UnimplementedChitChatServer
	savedMessages []proto.Message
}

func main() {
	server := &ChitChatServer{}

	server.startServer()
}

func (s *ChitChatServer) MessageToServer(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {
	//log.Println(msg)

	return &proto.Empty{}, nil
}

func (s *ChitChatServer) startServer() {
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

func (s *ChitChatServer) GetStream(pb *proto.Empty, stream proto.ChitChat_GetStreamServer) (err error) {

	// Based on the Message object's lamport timestamp, determine which message to stream out first

	for {
		err := stream.Send(&proto.Message{Msg: "Hello"})
		if err != nil {
			return err
		}
	}
	return nil
}
