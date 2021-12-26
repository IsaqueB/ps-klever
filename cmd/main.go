package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
)

func main() {
	//Setup and Run gRPC server
	client := database.NewMongoClient()
	if err := client.Connect(); err != nil {
		log.Fatalf("Error connecting to database. Error: %v", err)
	}
	defer client.Disconnect()

	s := rpc.NewGrpcServer(&client)
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":9000"
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting listener to grpc server. Error: %v", err)
	}
	proto.RegisterVoteServer(grpcServer, s)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error serving to grpc server. Error: %v", err)
	}
	fmt.Printf("Successfully started gRFC server in %v", port)
}
