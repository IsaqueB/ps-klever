package rpc

import (
	"context"
	"log"
	"strconv"

	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/IsaqueB/ps-klever/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

type Server interface {
	Insert(ctx context.Context, message *proto.InsertRequest) (*proto.InsertResponse, error)
	Get(ctx context.Context, message *proto.GetRequest) (*proto.GetResponse, error)
	UpdateOne(ctx context.Context, message *proto.UpdateOneRequest) (*proto.UpdateOneResponse, error)
	DeleteOne(ctx context.Context, message *proto.DeleteOneRequest) (*proto.DeleteOneResponse, error)
	ListVotesInVideo(req *proto.ListVotesInVideoRequest, stream proto.Vote_ListVotesInVideoServer) error
	ListVotesOfUser(req *proto.ListVotesOfUserRequest, stream proto.Vote_ListVotesOfUserServer) error
}

type server struct {
	client *database.MongoClient
}

// Create a new struct and sets it's client to the one in the function params
func CreateNewGrpcServer(client *database.MongoClient) Server {
	grpcServer := server{}
	grpcServer.setClient(client)
	return &grpcServer
}

// Set client to server in order to keep the same connection during operation
func (s *server) setClient(client *database.MongoClient) {
	s.client = client
}

// Create a new Vote from an USER to a VIDEO testar com o struct do protobuf
func (s *server) Insert(ctx context.Context, req *proto.InsertRequest) (*proto.InsertResponse, error) {
	log.Printf("INSERT VOTE - Recieved")
	// selecting collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// converting strings from request to objectId
	videoId, err := primitive.ObjectIDFromHex(req.Vote.Video)
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(req.Vote.User)
	if err != nil {
		return nil, err
	}
	// creating new document
	insertResult, err := voteCollection.InsertOne(ctx, database.VoteModel{
		ID:     primitive.NewObjectID(),
		Video:  videoId,
		User:   userId,
		Upvote: req.Vote.Upvote,
	})
	if err != nil {
		return nil, err
	}
	return &proto.InsertResponse{Id: insertResult.InsertedID.(primitive.ObjectID).Hex()}, nil
}

// Returns a Vote from an USER to a VIDEO
func (s *server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	log.Printf("GET UPVOTE - Recieved - ID TO BE QUERIED: %s", req.Id)
	// select collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// convert string from request to objectId
	voteId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	// query for the document
	var voteFound bson.M
	if err = voteCollection.FindOne(ctx, bson.M{"_id": voteId}).Decode(&voteFound); err != nil {
		return nil, err
	}
	// send message
	return &proto.GetResponse{Vote: &proto.VoteStruct{
		Id:     voteFound["_id"].(primitive.ObjectID).Hex(),
		Video:  voteFound["video"].(primitive.ObjectID).Hex(),
		User:   voteFound["user"].(primitive.ObjectID).Hex(),
		Upvote: voteFound["upvote"].(bool),
	}}, nil
}

// Modify vote's UPVOTE value which indicates if it is an UPVOTE or a DOWNVOTE
func (s *server) UpdateOne(ctx context.Context, req *proto.UpdateOneRequest) (*proto.UpdateOneResponse, error) {
	log.Printf("UPDATE UPVOTE - Recieved - ID: %s - CHANGE VALUE TO: %v", req.Id, req.NewValue)
	// select collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// convert string from request to objectId
	voteId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	// update document using the id and new upvote value got from request
	updateResult, err := voteCollection.UpdateByID(ctx, voteId, bson.M{"$set": bson.M{"upvote": req.NewValue}})
	if err != nil {
		return nil, err
	}
	// check to inform with the ID given does not correspond to a document in the database
	if updateResult.MatchedCount == 0 {
		return nil, status.Errorf(5, "Could not find the vote requested")
	}
	// send message
	return &proto.UpdateOneResponse{UpdateResult: "FOUND" +
		strconv.FormatInt(updateResult.MatchedCount, 10) +
		"MODIFIED:" +
		strconv.FormatInt(updateResult.ModifiedCount, 10)}, nil
}

// Remove an USER's vote to a VIDEO
func (s *server) DeleteOne(ctx context.Context, req *proto.DeleteOneRequest) (*proto.DeleteOneResponse, error) {
	log.Printf("DELETE UPVOTE - Recieved message from client: %s", req.Id)
	// select collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// convert string from request to objectId
	voteId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	deleteResult, err := voteCollection.DeleteOne(ctx, bson.M{"_id": voteId})
	if err != nil {
		return nil, err
	}
	if deleteResult.DeletedCount == 0 {
		return nil, status.Errorf(5, "Could not find the vote requested")
	}
	return &proto.DeleteOneResponse{DeleteResult: "DELETED:" + strconv.FormatInt(deleteResult.DeletedCount, 10)}, nil
}

// Queries all votes to a VIDEO and the UPVOTE count. Negative results to VoteCount means a video is more downvoted than upvoted
func (s *server) ListVotesInVideo(req *proto.ListVotesInVideoRequest, stream proto.Vote_ListVotesInVideoServer) error {
	log.Printf("GET VOTES OF VIDEO - Recieved - VIDEO TO BE QUERIED: %s", req.Id)
	ctx := context.Background()
	// selecting collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// converting string from request to objectId
	videoId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return err
	}
	// querying for votes of requested video
	cursor, err := voteCollection.Find(ctx, bson.M{"video": videoId})
	if err != nil {
		return err
	}
	// iterate cursor
	for cursor.Next(ctx) {
		var current bson.M
		// decode and parse info
		if err = cursor.Decode(&current); err != nil {
			return err
		}
		vote := &proto.VoteStruct{
			Id:     current["_id"].(primitive.ObjectID).Hex(),
			Video:  current["video"].(primitive.ObjectID).Hex(),
			User:   current["user"].(primitive.ObjectID).Hex(),
			Upvote: current["upvote"].(bool),
		}
		// send document
		err = stream.Send(&proto.ListVotesInVideoResponse{Vote: vote})
		if err != nil {
			return err
		}
	}
	return nil
}

//List all votes an USER made
func (s *server) ListVotesOfUser(req *proto.ListVotesOfUserRequest, stream proto.Vote_ListVotesOfUserServer) error {
	ctx := context.Background()
	// selecting collection
	voteCollection := (*s.client).GetClient().Database("ps-klever").Collection("vote")
	// converting string from request to objectId
	userId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return err
	}
	// querying for votes of requested video
	cursor, err := voteCollection.Find(ctx, bson.M{"user": userId})
	if err != nil {
		return err
	}
	// iterate cursor
	for cursor.Next(ctx) {
		var current bson.M
		// decode and parse info
		if err = cursor.Decode(&current); err != nil {
			return err
		}
		vote := &proto.VoteStruct{
			Id:     current["_id"].(primitive.ObjectID).Hex(),
			Video:  current["video"].(primitive.ObjectID).Hex(),
			User:   current["user"].(primitive.ObjectID).Hex(),
			Upvote: current["upvote"].(bool),
		}
		// send document
		err = stream.Send(&proto.ListVotesOfUserResponse{Vote: vote})
		if err != nil {
			return err
		}
	}
	return nil
}
