package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/DarkLordOfDeadstiny/Distributed_Mutual_Exclusion/gRPC"
	"google.golang.org/grpc"
)

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
	go joinServer(ctx, client)
	// exitChannel(ctx, client)


	//some code for getting access here, use for and gorutines?
	for {
		//keeps it running for now
		time.Sleep(time.Second * 5)// so it dont do the eat my cpu
	}
}

func joinServer(ctx context.Context, client gRPC.MessageServiceClient) {

	joinreq := gRPC.JoinRequest{SendersName: *sendername}
	responce, err := client.Join(ctx, &joinreq)
	if err != nil {
		log.Fatalf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := responce.Recv()
			if err == io.EOF {
				fmt.Println("error")
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to Join server. \nErr: %v", err)
			}

			if in.LamportTime > *lamportTime {
				*lamportTime = in.LamportTime + 1
			} else {
				*lamportTime++
			}

			//some resquest something thing here
			time.Sleep(time.Duration(rand.Intn(2))) //sleep for a random time up to 2 sec hilsen patrick
		}
	}()

	<-waitc
}

//add methods for requesting things