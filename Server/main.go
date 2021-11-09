package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/DarkLordOfDeadstiny/Distributed_Mutual_Exclusion/gRPC"

	"google.golang.org/grpc"
)

var requestQueue = make(chan int, 10)

type Server struct {
	gRPC.UnimplementedMessageServiceServer
	messageChan []chan *gRPC.JoinRequest
}

func main() {
	fmt.Println(".:server is starting:.")

	// Create listener tcp on port 5400
	list, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("Failed to listen on port 5400: %v", err)
	}

	//Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	//connect to log file
	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gRPC.RegisterMessageServiceServer(grpcServer, newServer())
	grpcServer.Serve(list)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func newServer() *Server {
	s := &Server{}
	fmt.Println(s)
	return s
}

func (s *Server) Join(ch *gRPC.JoinRequest, msgStream gRPC.MessageService_JoinServer) error {
	
	for {
		select {
		case <-msgStream.Context().Done():
			//what happens when someone leaves
			return nil
			// case something:
	}


}

