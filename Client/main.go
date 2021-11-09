package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	// "strings"

	"github.com/DarkLordOfDeadstiny/Distributed_Mutual_Exclusion/gRPC"
	"google.golang.org/grpc"
)

var channelname = flag.String("Channel", "default", "Chitty-Chat")
var sendername = flag.String("sender", "default", "Senders name")
var tcpServer = flag.String("server", "localhost:5400", "Tcp server")
var lamportTime = flag.Int64("time", 0, "lamportTimeStamp")

func main() {
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Fail to Dial : %v", err)
	}

	defer conn.Close()

	ctx := context.Background()
	client := gRPC.NewMessageServiceClient(conn)

	fmt.Println("--- join channel ---")
	go joinChannel(ctx, client)
	// exitChannel(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text())
	}

}

func joinChannel(ctx context.Context, client gRPC.MessageServiceClient) {

	joinreq := gRPC.JoinRequest{ChanName: *channelname, SendersName: *sendername}
	stream, err := client.Join(ctx, &joinreq)
	if err != nil {
		log.Fatalf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}

	fmt.Printf("Joined channel: %v \n", *channelname)

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}

			if *sendername != in.Sender {
				if in.LamportTime > *lamportTime {
					*lamportTime = in.LamportTime + 1
				} else {
					*lamportTime++
				}
				fmt.Printf("MESSAGE: (%v) %v: %v \n", *lamportTime, in.Sender, in.Message)
			}
		}
	}()

	<-waitc
}
