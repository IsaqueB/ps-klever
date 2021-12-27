package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	pb "github.com/IsaqueB/ps-klever/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	errors := make(chan error)
	//Setup and Run HTTP Server
	go func() {
		mux := runtime.NewServeMux()
		client := database.NewMongoClient()
		if err := client.Connect(); err != nil {
			log.Fatalf("Error connecting to database. Error: %v", err)
		}
		defer client.Disconnect()

		s := rpc.NewGrpcServer(&client)
		grpcServer := grpc.NewServer()
		defer grpcServer.GracefulStop()
		pb.RegisterVoteHandlerServer(context.Background(), mux, s)
		port := ":" + os.Getenv("PORT")
		if port == ":" {
			port = ":9000"
		}
		errors <- http.ListenAndServe(port, mux)
	}()
	//Setup and Run gRPC Server
	// Not used since heroku doesn't support HTTP/2
	go func() {
		client := database.NewMongoClient()
		if err := client.Connect(); err != nil {
			log.Fatalf("Error connecting to database. Error: %v", err)
		}
		defer client.Disconnect()

		s := rpc.NewGrpcServer(&client)
		grpcServer := grpc.NewServer()
		defer grpcServer.GracefulStop()

		port := ":" + os.Getenv("PORT_GRPC")
		if port == ":" {
			port = ":9001"
		}
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("Error starting listener to grpc server. Error: %v", err)
		}

		pb.RegisterVoteServer(grpcServer, s)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Error serving to grpc server. Error: %v", err)
		}
		fmt.Printf("Successfully started gRFC server in %v", port)
	}()
	for err := range errors {
		log.Fatal(err)
		return
	}
}
