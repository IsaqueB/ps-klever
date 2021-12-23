package main

import (
	"log"
	"net"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Failed to listen. Error: %v", err)
	}

	client := database.CreateNewMongoClient()
	if err = client.Connect(); err != nil {
		log.Fatalf("Error connecting to database. Error: %v", err)
	}

	defer client.Disconnect()

	if err != nil {
		log.Fatalf("Failed to listen. Error: %v", err)
	}

	s := rpc.CreateNewGrpcServer(&client)

	grpcServer := grpc.NewServer()

	proto.RegisterVoteServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server. %v", err)
	}

	defer grpcServer.GracefulStop()
}
