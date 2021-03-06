package rpc_test

import (
	"context"
	"testing"

	"github.com/IsaqueB/ps-klever/cmd/rpc"
	"github.com/IsaqueB/ps-klever/pkg/database"
	pb "github.com/IsaqueB/ps-klever/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewGrpcServer(t *testing.T) {
	client := database.NewMongoClient()
	s := rpc.NewGrpcServer(&client)
	assert.Equal(t, *s.GetClient(), client, "server's client and the created should be the same")
}

func initAServer() (rpc.Server, error) {
	client := database.NewMongoClient()
	if err := client.Connect(); err != nil {
		return nil, err
	}
	s := rpc.NewGrpcServer(&client)
	s.SetDatabase(rpc.TEST_DB)
	return s, nil
}

func TestInsert(t *testing.T) {
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	mock_req := &pb.InsertRequest{
		Vote: &pb.VoteStruct{
			Video:  mock_id,
			User:   mock_id,
			Upvote: true,
		},
	}
	s, err := initAServer()
	defer (*s.GetClient()).Disconnect()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	res, err := s.Insert(mock_ctx, mock_req)
	if err != nil {
		t.Errorf("Error inside Insert: %v", err)
	}
	if objectId, err := primitive.ObjectIDFromHex(res.Id); err != nil {
		t.Errorf("Did not return an id, %v. %v", objectId, err)
	}
}

func TestGet(t *testing.T) {
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	vote := pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	mock_req := &pb.InsertRequest{
		Vote: &vote,
	}
	s, err := initAServer()
	defer (*s.GetClient()).Disconnect()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	res_insert, err := s.Insert(mock_ctx, mock_req)
	if err != nil {
		t.Errorf("Error inside Insert: %v", err)
	}
	vote.Id = res_insert.Id
	res_get, err := s.Get(mock_ctx, &pb.GetRequest{
		Id: res_insert.Id,
	})
	if err != nil {
		t.Errorf("Error inside Get: %v", err)
	}
	assert.Equal(t, &vote, res_get.GetVote())
}
func TestUpdateOne(t *testing.T) {
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	vote := pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	mock_req := &pb.InsertRequest{
		Vote: &vote,
	}
	s, err := initAServer()
	defer (*s.GetClient()).Disconnect()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	res_insert, err := s.Insert(mock_ctx, mock_req)
	if err != nil {
		t.Errorf("Error inside Insert: %v", err)
	}
	res_updateOne, err := s.UpdateOne(mock_ctx, &pb.UpdateOneRequest{
		Id:       res_insert.Id,
		NewValue: false,
	})
	if err != nil {
		t.Errorf("Error inside Update: %v", err)
	}
	assert.Equal(t, res_updateOne.GetMatched(), int32(1), "The amount of documents matched should be one")
	assert.Equal(t, res_updateOne.GetModified(), int32(1), "The amount of documents modified should be one")
}
func TestDeleteOne(t *testing.T) {
	mock_ctx := context.Background()
	mock_id := primitive.NewObjectID().Hex()
	vote := pb.VoteStruct{
		Video:  mock_id,
		User:   mock_id,
		Upvote: true,
	}
	mock_req := &pb.InsertRequest{
		Vote: &vote,
	}
	s, err := initAServer()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	defer (*s.GetClient()).Disconnect()
	res_insert, err := s.Insert(mock_ctx, mock_req)
	if err != nil {
		t.Errorf("Error inside Insert: %v", err)
	}
	res_delete, err := s.DeleteOne(mock_ctx, &pb.DeleteOneRequest{
		Id: res_insert.Id,
	})
	if err != nil {
		t.Errorf("Error inside Delete: %v", err)
	}
	assert.Equal(t, res_delete.GetDeleted(), int32(1), "The amount of documents deleted should be one")
}

func TestListVotesInVideo(t *testing.T) {
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
	// Setup server
	s, err := initAServer()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	defer (*s.GetClient()).Disconnect()
	// Populate database
	res_0, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_0})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_0.Id = res_0.GetId()
	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_1})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_2})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	res_3, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_3})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_3.Id = res_3.GetId()
	//Get info from database
	res, err := s.ListVotesInVideo(context.Background(), &pb.ListVotesInVideoRequest{Id: mock_id_0})
	if err != nil {
		t.Errorf("Error in ListVotesInVideo. %v", err)
	}
	assert.Equal(t, []*pb.VoteStruct{&mock_Vote_0, &mock_Vote_3}, res.Vote)
}

func ListVotesOfUser(t *testing.T) {
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
		Video:  mock_id_2,
		User:   mock_id_1,
		Upvote: true,
	}
	mock_Vote_2 := pb.VoteStruct{
		Video:  mock_id_0,
		User:   mock_id_2,
		Upvote: true,
	}
	mock_Vote_3 := pb.VoteStruct{
		Video:  mock_id_1,
		User:   mock_id_0,
		Upvote: true,
	}
	// Setup server
	s, err := initAServer()
	if err != nil {
		t.Errorf("Error setting up server. %v", err)
	}
	defer (*s.GetClient()).Disconnect()
	// Populate database
	res_0, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_0})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_0.Id = res_0.GetId()
	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_1})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_2})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	res_3, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_3})
	if err != nil {
		t.Errorf("Error in Insert. %v", err)
	}
	mock_Vote_3.Id = res_3.GetId()
	//Get info from database
	res, err := s.ListVotesOfUser(context.Background(), &pb.ListVotesOfUserRequest{Id: mock_id_0})
	if err != nil {
		t.Errorf("Error in ListVotesOfUser. %v", err)
	}
	assert.Equal(t, []*pb.VoteStruct{&mock_Vote_0, &mock_Vote_3}, res.Vote)
}

// These were used when protobuf was returning a stream. but since the http handler can't handle them yet
// I've changed them to returning an array
// START OF TestListVotesInVideo STRAEM TEST
// type mockVoteListVotesInVideoServer struct {
// 	grpc.ServerStream
// 	Results []*pb.VoteStruct
// }

// func (x *mockVoteListVotesInVideoServer) Send(m *pb.ListVotesInVideoResponse) error {
// 	x.Results = append(x.Results, m.Vote)
// 	return nil
// }

// func TestListVotesInVideo(t *testing.T) {
// 	//Create mocks
// 	mock_id_0 := primitive.NewObjectID().Hex()
// 	mock_id_1 := primitive.NewObjectID().Hex()
// 	mock_id_2 := primitive.NewObjectID().Hex()
// 	mock_Vote_0 := pb.VoteStruct{
// 		Video:  mock_id_0,
// 		User:   mock_id_2,
// 		Upvote: true,
// 	}
// 	mock_Vote_1 := pb.VoteStruct{
// 		Video:  mock_id_1,
// 		User:   mock_id_2,
// 		Upvote: true,
// 	}
// 	mock_Vote_2 := pb.VoteStruct{
// 		Video:  mock_id_1,
// 		User:   mock_id_2,
// 		Upvote: true,
// 	}
// 	mock_Vote_3 := pb.VoteStruct{
// 		Video:  mock_id_0,
// 		User:   mock_id_2,
// 		Upvote: true,
// 	}
// 	//Setup server
// 	s, err := initAServer()
// 	if err != nil {
// 		t.Errorf("Error setting up server. %v", err)
// 	}
// 	defer (*s.GetClient()).Disconnect()
// 	//Populate database
// 	res_0, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_0})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	mock_Vote_0.Id = res_0.GetId()
// 	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_1})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_2})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	res_3, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_3})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	mock_Vote_3.Id = res_3.GetId()
// 	streamMock := mockVoteListVotesInVideoServer{}
// 	err = s.ListVotesInVideo(&pb.ListVotesInVideoRequest{Id: mock_id_0}, &streamMock)
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	// checks for first vote
// 	assert.Equal(t, mock_Vote_0.Id, streamMock.Results[0].GetId(), "Ids should be equal")
// 	assert.Equal(t, mock_Vote_0.Video, streamMock.Results[0].GetVideo(), "Video's id should be equal")
// 	assert.Equal(t, mock_Vote_0.User, streamMock.Results[0].GetUser(), "User's id should be equal")
// 	assert.Equal(t, mock_Vote_0.Upvote, streamMock.Results[0].GetUpvote(), "Upvote's value should be equal")
// 	// checks for second vote
// 	assert.Equal(t, mock_Vote_3.Id, streamMock.Results[1].GetId(), "Ids should be equal")
// 	assert.Equal(t, mock_Vote_3.Video, streamMock.Results[1].GetVideo(), "Video's id should be equal")
// 	assert.Equal(t, mock_Vote_3.User, streamMock.Results[1].GetUser(), "User's id should be equal")
// 	assert.Equal(t, mock_Vote_3.Upvote, streamMock.Results[1].GetUpvote(), "Upvote's value should be equal")
// }

// END OF TestListVotesInVideo's STREAM TEST

// // START OF TestListVotesOfUser STRAEM TEST
// type mockVoteListVotesOfUserServer struct {
// 	grpc.ServerStream
// 	Results []*pb.VoteStruct
// }

// func (x *mockVoteListVotesOfUserServer) Send(m *pb.ListVotesOfUserResponse) error {
// 	x.Results = append(x.Results, m.Vote)
// 	return nil
// }

// func TestListVotesOfUser(t *testing.T) {
// 	//Create mocks
// 	mock_id_0 := primitive.NewObjectID().Hex()
// 	mock_id_1 := primitive.NewObjectID().Hex()
// 	mock_id_2 := primitive.NewObjectID().Hex()
// 	mock_Vote_0 := pb.VoteStruct{
// 		Video:  mock_id_2,
// 		User:   mock_id_1,
// 		Upvote: true,
// 	}
// 	mock_Vote_1 := pb.VoteStruct{
// 		Video:  mock_id_2,
// 		User:   mock_id_0,
// 		Upvote: true,
// 	}
// 	mock_Vote_2 := pb.VoteStruct{
// 		Video:  mock_id_2,
// 		User:   mock_id_0,
// 		Upvote: true,
// 	}
// 	mock_Vote_3 := pb.VoteStruct{
// 		Video:  mock_id_2,
// 		User:   mock_id_1,
// 		Upvote: true,
// 	}
// 	//Setup server
// 	s, err := initAServer()
// 	if err != nil {
// 		t.Errorf("Error setting up server. %v", err)
// 	}
// 	defer (*s.GetClient()).Disconnect()
// 	//Populate database
// 	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_0})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	res_1, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_1})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	mock_Vote_1.Id = res_1.GetId()
// 	res_2, err := s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_2})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	mock_Vote_2.Id = res_2.GetId()
// 	_, err = s.Insert(context.Background(), &pb.InsertRequest{Vote: &mock_Vote_3})
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	streamMock := mockVoteListVotesOfUserServer{}
// 	err = s.ListVotesOfUser(&pb.ListVotesOfUserRequest{Id: mock_id_0}, &streamMock)
// 	if err != nil {
// 		t.Errorf("Error in Insert. %v", err)
// 	}
// 	// checks of first vote
// 	assert.Equal(t, mock_Vote_1.Id, streamMock.Results[0].GetId(), "Ids should be equal")
// 	assert.Equal(t, mock_Vote_1.Video, streamMock.Results[0].GetVideo(), "Video's id should be equal")
// 	assert.Equal(t, mock_Vote_1.User, streamMock.Results[0].GetUser(), "User's id should be equal")
// 	assert.Equal(t, mock_Vote_1.Upvote, streamMock.Results[0].GetUpvote(), "Upvote's value should be equal")
// 	// checks of second vote
// 	assert.Equal(t, mock_Vote_2.Id, streamMock.Results[1].GetId(), "Ids should be equal")
// 	assert.Equal(t, mock_Vote_2.Video, streamMock.Results[1].GetVideo(), "Video's id should be equal")
// 	assert.Equal(t, mock_Vote_2.User, streamMock.Results[1].GetUser(), "User's id should be equal")
// 	assert.Equal(t, mock_Vote_2.Upvote, streamMock.Results[1].GetUpvote(), "Upvote's value should be equal")
// }
// END OF TestListVotesOfUser's STREAM TEST
