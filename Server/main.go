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

type Server struct {
	gRPC.UnimplementedMessageServiceServer
	Clients map[string]*gRPC.JoinRequest
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
	s := &Server{Clients: make(map[string]*gRPC.JoinRequest)}
	fmt.Println(s) //prints the server struct to console
	return s
}

func (s *Server) Join(joinRequest *gRPC.JoinRequest, msgStream gRPC.MessageService_JoinServer) error {
	s.Clients[joinRequest.SendersName] = joinRequest //adds new client to the list

	log.Printf("%s Joined the server", joinRequest.SendersName)

	for {
		select {
			case <-msgStream.Context().Done():
				//what happens when someone leaves
				return nil
				// case 'something that should run all the time':
		}
	}
}


func (s *Server) Leave(ctx context.Context, leaveRequest *gRPC.LeaveRequest) (*gRPC.LeaveResponse, error) {
	delete(s.Clients, leaveRequest.SendersName)
	return &gRPC.LeaveResponse{Status: "left"}, nil
}

func Entry(ctx context.Context, request *gRPC.EntryRequest) (*gRPC.EntryResponse, error) {
	select{
		case requestQueue <- request:
		default:
			return &gRPC.EntryResponse{Status: "queue is full, try again later"},
					errors.New("action not posible")
	}
	

	return &gRPC.EntryResponse{}, nil
}

func ResourceAccess(ctx context.Context, request *gRPC.AccessRequest) (*gRPC.AccessResponse, error){
	return nil, errors.New("not implemented")
}

func Exit(ctx context.Context, request *gRPC.ExitRequest) (*gRPC.ExitResponse, error){
	return nil, errors.New("not implemented")
}