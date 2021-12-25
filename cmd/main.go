package main

import (
	"context"
	"log"
	"net"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/IsaqueB/ps-klever/pkg/http_srv"
	"github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
)

func main() {
	errors := make(chan error)
	//Setup and Run gRPC server
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen. Error: %v", err)
	}
	client := database.NewMongoClient()
	if err = client.Connect(); err != nil {
		log.Fatalf("Error connecting to database. Error: %v", err)
	}
	defer client.Disconnect()
	if err != nil {
		log.Fatalf("Failed to listen. Error: %v", err)
	}
	s := rpc.NewGrpcServer(&client)
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()
	proto.RegisterVoteServer(grpcServer, s)
	go func() {
		errors <- grpcServer.Serve(lis)
	}()
	defer grpcServer.GracefulStop()
	//Setup and Run HTTP server
	http_s, err := http_srv.NewHTTPServer()
	if err != nil {
		log.Fatalf("Error creating and connecting the HTTP server. %v", err)
	}
	go func() {
		errors <- http_s.GetServer().ListenAndServe()
	}()
	defer http_s.GetServer().Shutdown(context.Background())
	for err := range errors {
		log.Fatalf("%v", err)
		return
	}
}
