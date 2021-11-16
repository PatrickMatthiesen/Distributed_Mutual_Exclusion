package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/DarkLordOfDeadstiny/Distributed_Mutual_Exclusion/gRPC"

	"google.golang.org/grpc"
)

var requestQueue = make(chan *gRPC.EntryRequest, 10)
var usedBy *gRPC.JoinRequest

type Server struct {
	gRPC.UnimplementedMessageServiceServer
	Clients map[string]*gRPC.JoinRequest
	turns   map[string]chan bool
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
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	s := &Server{
		Clients: make(map[string]*gRPC.JoinRequest),
		turns:   make(map[string]chan bool),
	}
	fmt.Println(s) //prints the server struct to console
	return s
}

func (s *Server) Join(joinRequest *gRPC.JoinRequest, msgStream gRPC.MessageService_JoinServer) error {
	s.Clients[joinRequest.SendersName] = joinRequest //adds new client to the list
	s.turns[joinRequest.SendersName] = make(chan bool, 1)

	log.Printf("%s Joined the server", joinRequest.SendersName)

	return nil
}

func (s *Server) Leave(ctx context.Context, leaveRequest *gRPC.LeaveRequest) (*gRPC.LeaveResponse, error) {
	delete(s.Clients, leaveRequest.SendersName)
	return &gRPC.LeaveResponse{Status: "left"}, nil
}

func (s *Server) Entry(ctx context.Context, request *gRPC.EntryRequest) (*gRPC.EntryResponse, error) {
	log.Printf("Client %s attempts to retrieve the access token \n", request.SendersName)
	select {
	case requestQueue <- request:
		if usedBy == nil {
			usedBy = s.Clients[request.SendersName]
			s.giveAccessToNext()
		}
		log.Printf("client %s waits for its turn\n", request.SendersName)
		<-s.turns[request.SendersName]
		return &gRPC.EntryResponse{Status: "200"}, nil
	default:
		return &gRPC.EntryResponse{Status: "409"},
			errors.New("action not posible, queue is full, try again later")
	}
}

func (s *Server) ResourceAccess(ctx context.Context, request *gRPC.AccessRequest) (*gRPC.AccessResponse, error) {
	usedBy = s.Clients[request.SendersName]
	fmt.Printf("Client %s has accessed the resource \n", request.SendersName)
	log.Printf("Client %s has accessed the resource \n", request.SendersName)

	return &gRPC.AccessResponse{Status: "200"}, nil
}

func (s *Server) Exit(ctx context.Context, request *gRPC.ExitRequest) (*gRPC.ExitResponse, error) {
	log.Printf("Client %s exits the resource \n", request.SendersName)
	s.giveAccessToNext()

	return &gRPC.ExitResponse{Status: "200"}, nil
}

func (s *Server) giveAccessToNext() {
	select {
	case request := <-requestQueue:
		fmt.Printf("Access to resource has been given to client: (%s)\n", request.SendersName)
		log.Printf("Access to resource has been given to client: (%s)\n", request.SendersName)
		s.turns[request.SendersName] <- true
	default:
		usedBy = nil
	}
}
