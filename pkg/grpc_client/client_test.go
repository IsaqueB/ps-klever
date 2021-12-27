package grpc_client_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/IsaqueB/ps-klever/pkg/grpc_client"
	pb "github.com/IsaqueB/ps-klever/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

func init() {
	// start server
	client := database.NewMongoClient()
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	s := rpc.NewGrpcServer(&client)
	s.SetDatabase(rpc.TEST_DB)
	grpcServer := grpc.NewServer()

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":9000"
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterVoteServer(grpcServer, s)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with %v", err)
		}
	}()
}

func initClient() (grpc_client.Client, error) {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":9000"
	}
	client, err := grpc_client.NewGrpcClient(port)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func compareVotes(t *testing.T, a *pb.VoteStruct, b *pb.VoteStruct) {
	assert.Equal(t, a.GetId(), b.GetId())
	assert.Equal(t, a.GetVideo(), b.GetVideo())
	assert.Equal(t, a.GetUser(), b.GetUser())
	assert.Equal(t, a.GetUpvote(), b.GetUpvote())
}

func TestNewGrpcClient(t *testing.T) {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":9000"
	}
	if _, err := grpc_client.NewGrpcClient(port); err != nil {
		t.Fatalf("Failed to generate new grpc client and listen to %v. %v", port, err)
	}
}

func TestInsert(t *testing.T) {
	c, err := initClient()
	defer c.Disconnect()
	if err != nil {
		t.Fatalf("Error initializing client. %v", err)
	}
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	mock_req := &pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	_, err = c.Insert(mock_ctx, mock_req)
	if err != nil {
		t.Fatalf("Error in Insert. %v", err)
	}
}

func TestGet(t *testing.T) {
	c, err := initClient()
	if err != nil {
		t.Fatalf("Error creating client. %v", err)
	}
	defer c.Disconnect()
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	mock_req := &pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	inserted_id, _ := c.Insert(mock_ctx, mock_req)
	mock_req.Id = inserted_id
	vote, err := c.Get(mock_ctx, inserted_id)
	if err != nil {
		t.Fatalf("Error in Get. %v", vote)
	}
	compareVotes(t, mock_req, vote)
}

func TestUpdateOne(t *testing.T) {
	c, err := initClient()
	if err != nil {
		t.Fatalf("Error creating client. %v", err)
	}
	defer c.Disconnect()
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	mock_req := &pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	inserted_id, _ := c.Insert(mock_ctx, mock_req)
	matched, modifeid, err := c.UpdateOne(mock_ctx, inserted_id, false)
	if err != nil {
		t.Errorf("Error inside Update: %v", err)
	}
	assert.Equal(t, matched, int32(1), "The amount of documents matched should be one")
	assert.Equal(t, modifeid, int32(1), "The amount of documents modified should be one")
}

func TestDeleteOne(t *testing.T) {
	c, err := initClient()
	if err != nil {
		t.Fatalf("Error creating client. %v", err)
	}
	defer c.Disconnect()
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	mock_req := &pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	inserted_id, _ := c.Insert(mock_ctx, mock_req)
	deleted, err := c.DeleteOne(mock_ctx, inserted_id)
	if err != nil {
		t.Errorf("Error inside Update: %v", err)
	}
	assert.Equal(t, deleted, int32(1), "The amount of documents deleted should be one")
}

func TestListVotesInVideo(t *testing.T) {
	c, err := initClient()
	if err != nil {
		t.Fatalf("Error creating client. %v", err)
	}
	defer c.Disconnect()
	// Create mocks
	mock_id_0 := primitive.NewObjectID().Hex()
	mock_id_1 := primitive.NewObjectID().Hex()
	mock_id_2 := primitive.NewObjectID().Hex()
	mock_Vote_0 := pb.VoteStruct{
		Video:  mock_id_0,
		User:   mock_id_2,
		Upvote: true,
	}
	mock_Vote_1 := pb.VoteStruct{
		Video:  mock_id_1,
		User:   mock_id_2,
		Upvote: true,
	}
	mock_Vote_2 := pb.VoteStruct{
		Video:  mock_id_1,
		User:   mock_id_2,
		Upvote: true,
	}
	mock_Vote_3 := pb.VoteStruct{
		Video:  mock_id_0,
		User:   mock_id_2,
		Upvote: true,
	}
	// Populate database
	id0, err := c.Insert(context.Background(), &mock_Vote_0)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_0.Id = id0
	_, err = c.Insert(context.Background(), &mock_Vote_1)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	_, err = c.Insert(context.Background(), &mock_Vote_2)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	id3, err := c.Insert(context.Background(), &mock_Vote_3)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_3.Id = id3
	// Make request
	response, err := c.GetClient().ListVotesInVideo(context.Background(), &pb.ListVotesInVideoRequest{
		Id: mock_id_0,
	})
	if err != nil {
		t.Fatalf("Error in ListVotesInVideo. %v", err)
	}
	votes := []*pb.VoteStruct{&mock_Vote_0, &mock_Vote_3}
	for i, vote := range response.Vote {
		compareVotes(t, votes[i], vote)
	}
}

func TestListVotesOfUser(t *testing.T) {
	c, err := initClient()
	if err != nil {
		t.Fatalf("Error creating client. %v", err)
	}
	defer c.Disconnect()
	// Create mocks
	mock_id_0 := primitive.NewObjectID().Hex()
	mock_id_1 := primitive.NewObjectID().Hex()
	mock_id_2 := primitive.NewObjectID().Hex()
	mock_Vote_0 := pb.VoteStruct{
		Video:  mock_id_2,
		User:   mock_id_0,
		Upvote: true,
	}
	mock_Vote_1 := pb.VoteStruct{
		Video:  mock_id_1,
		User:   mock_id_0,
		Upvote: true,
	}
	mock_Vote_2 := pb.VoteStruct{
		Video:  mock_id_2,
		User:   mock_id_1,
		Upvote: true,
	}
	mock_Vote_3 := pb.VoteStruct{
		Video:  mock_id_0,
		User:   mock_id_2,
		Upvote: true,
	}
	// Populate database
	id0, err := c.Insert(context.Background(), &mock_Vote_0)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_0.Id = id0
	id1, err := c.Insert(context.Background(), &mock_Vote_1)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_1.Id = id1
	_, err = c.Insert(context.Background(), &mock_Vote_2)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	_, err = c.Insert(context.Background(), &mock_Vote_3)
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	// Make request
	response, err := c.GetClient().ListVotesOfUser(context.Background(), &pb.ListVotesOfUserRequest{
		Id: mock_id_0,
	})
	if err != nil {
		t.Fatalf("Error in ListVotesOfUser. %v", err)
	}
	votes := []*pb.VoteStruct{&mock_Vote_0, &mock_Vote_1}
	for i, vote := range response.Vote {
		compareVotes(t, votes[i], vote)
	}
}
