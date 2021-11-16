package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	// "math/rand"
	"os"
	"time"

	"github.com/DarkLordOfDeadstiny/Distributed_Mutual_Exclusion/gRPC"
	"google.golang.org/grpc"
)

var sendername = flag.String("sender", "default", "Senders name")
var tcpServer = flag.String("server", "localhost:5400", "Tcp server")
var lamportTime = flag.Int64("time", 0, "lamportTimeStamp")

var ctx context.Context
var client gRPC.MessageServiceClient

func main() {
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Fail to Dial : %v", err)
	}

	//connect to log file
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	defer conn.Close()

	ctx = context.Background()
	client = gRPC.NewMessageServiceClient(conn)

	fmt.Println("--- join channel ---")
	go joinServer()
	// exitChannel(ctx, client)
	sendEntryRequest()

}

func joinServer() {

	joinreq := gRPC.JoinRequest{SendersName: *sendername}
	response, err := client.Join(ctx, &joinreq)
	if err != nil {
		fmt.Printf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		_, err := response.Recv()
		if err == io.EOF {
			fmt.Println("error")
			close(waitc)
			return
		}
		if err != nil {
			fmt.Printf("Failed to Join server. \nErr: %v", err)
		}
	}()

	<-waitc
}

func leaveServer() {
	request := &gRPC.LeaveRequest{
		SendersName: *sendername,
	}
	client.Leave(ctx, request)
	tcpServer = nil
	fmt.Printf("client: %s left the server\n", *sendername)
}

func sendEntryRequest() {
	fmt.Printf("client: %s sends entry request\n", *sendername)
	request := &gRPC.EntryRequest{
		SendersName: *sendername,
		LamportTime: *lamportTime,
	}
	response, err := client.Entry(ctx, request)

	if err != nil {
		log.Fatalf("client.sendEntryRequest() throws: %v", err)
	}

	switch response.Status {
	case "200":
		getAccess()
		sendEntryRequest()
		break
	case "409":
		sendEntryRequest()

	}
}

func getAccess() {
	request := &gRPC.AccessRequest{
		SendersName: *sendername,
		LamportTime: *lamportTime,
	}
	response, err := client.ResourceAccess(ctx, request)

	if err != nil {
		log.Fatalf("client.getAccess(ctx, &channel) throws: %v", err)
	}

	if response != nil {

	}

	time.Sleep(time.Second * 1)
	returnAccess()
}

func returnAccess() {
	request := &gRPC.ExitRequest{
		SendersName: *sendername,
		LamportTime: *lamportTime,
	}
	client.Exit(ctx, request)
}
