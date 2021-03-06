package rpc

import (
	"context"
	"log"

	"github.com/IsaqueB/ps-klever/pkg/database"
	pb "github.com/IsaqueB/ps-klever/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

const (
	MAIN_DB int = iota
	TEST_DB int = iota
)

var (
	db_string = map[int]string{
		MAIN_DB: "ps-klever",
		TEST_DB: "ps-klever-test",
	}
)

type Server interface {
	Insert(ctx context.Context, message *pb.InsertRequest) (*pb.InsertResponse, error)
	Get(ctx context.Context, message *pb.GetRequest) (*pb.GetResponse, error)
	UpdateOne(ctx context.Context, message *pb.UpdateOneRequest) (*pb.UpdateOneResponse, error)
	DeleteOne(ctx context.Context, message *pb.DeleteOneRequest) (*pb.DeleteOneResponse, error)
	ListVotesInVideo(ctx context.Context, req *pb.ListVotesInVideoRequest) (*pb.ListVotesInVideoResponse, error)
	ListVotesOfUser(ctx context.Context, req *pb.ListVotesOfUserRequest) (*pb.ListVotesOfUserResponse, error)
	GetClient() *database.MongoClient
	SetDatabase(index int)
	pb.UnsafeVoteServer
}

type server struct {
	client   *database.MongoClient
	database string
	pb.UnimplementedVoteServer
}

// Create a new struct and sets it's client to the one in the function params
func NewGrpcServer(client *database.MongoClient) Server {
	grpcServer := server{}
	grpcServer.setClient(client)
	grpcServer.SetDatabase(MAIN_DB)
	return &grpcServer
}

// Set client to server in order to keep the same connection during operation

func (s *server) setClient(client *database.MongoClient) {
	s.client = client
}

func (s *server) GetClient() *database.MongoClient {
	return s.client
}

func (s *server) SetDatabase(index int) {
	s.database = db_string[index]
}

// Create a new Vote from an USER to a VIDEO testar com o struct do pbbuf
func (s *server) Insert(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	log.Println("INSERT VOTE - Recieved")
	// selecting collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
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
	return &pb.InsertResponse{Id: insertResult.InsertedID.(primitive.ObjectID).Hex()}, nil
}

// Returns a Vote from an USER to a VIDEO
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("GET UPVOTE - Recieved - ID TO BE QUERIED: %s", req.Id)
	// select collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
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
	return &pb.GetResponse{Vote: &pb.VoteStruct{
		Id:     voteFound["_id"].(primitive.ObjectID).Hex(),
		Video:  voteFound["video"].(primitive.ObjectID).Hex(),
		User:   voteFound["user"].(primitive.ObjectID).Hex(),
		Upvote: voteFound["upvote"].(bool),
	}}, nil
}

// Modify vote's UPVOTE value which indicates if it is an UPVOTE or a DOWNVOTE
func (s *server) UpdateOne(ctx context.Context, req *pb.UpdateOneRequest) (*pb.UpdateOneResponse, error) {
	log.Printf("UPDATE UPVOTE - Recieved - ID: %s - CHANGE VALUE TO: %v", req.Id, req.NewValue)
	// select collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
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
	return &pb.UpdateOneResponse{
		Matched:  int32(updateResult.MatchedCount),
		Modified: int32(updateResult.ModifiedCount),
	}, nil
}

// Remove an USER's vote to a VIDEO
func (s *server) DeleteOne(ctx context.Context, req *pb.DeleteOneRequest) (*pb.DeleteOneResponse, error) {
	log.Printf("DELETE UPVOTE - Recieved message from client: %s", req.Id)
	// select collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
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
	return &pb.DeleteOneResponse{
		Deleted: int32(deleteResult.DeletedCount),
	}, nil
}

// Queries all votes to a VIDEO and the UPVOTE count. Negative results to VoteCount means a video is more downvoted than upvoted
func (s *server) ListVotesInVideo(ctx context.Context, req *pb.ListVotesInVideoRequest) (*pb.ListVotesInVideoResponse, error) {
	log.Printf("GET VOTES OF VIDEO - Recieved - VIDEO TO BE QUERIED: %s", req.Id)
	// selecting collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
	// converting string from request to objectId
	videoId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	// querying for votes of requested video
	cursor, err := voteCollection.Find(ctx, bson.M{"video": videoId})
	if err != nil {
		return nil, err
	}
	// iterate cursor
	var votes []*pb.VoteStruct
	for cursor.Next(ctx) {
		var current bson.M
		// decode and parse info
		if err = cursor.Decode(&current); err != nil {
			return nil, err
		}
		vote := &pb.VoteStruct{
			Id:     current["_id"].(primitive.ObjectID).Hex(),
			Video:  current["video"].(primitive.ObjectID).Hex(),
			User:   current["user"].(primitive.ObjectID).Hex(),
			Upvote: current["upvote"].(bool),
		}
		// send document
		votes = append(votes, vote)
	}
	return &pb.ListVotesInVideoResponse{
		Vote: votes,
	}, nil
}

//List all votes an USER made
func (s *server) ListVotesOfUser(ctx context.Context, req *pb.ListVotesOfUserRequest) (*pb.ListVotesOfUserResponse, error) {
	// selecting collection
	voteCollection := (*s.client).GetClient().Database(s.database).Collection("vote")
	// converting string from request to objectId
	userId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	// querying for votes of requested video
	cursor, err := voteCollection.Find(ctx, bson.M{"user": userId})
	if err != nil {
		return nil, err
	}
	// iterate cursor
	var votes []*pb.VoteStruct
	for cursor.Next(ctx) {
		var current bson.M
		// decode and parse info
		if err = cursor.Decode(&current); err != nil {
			return nil, err
		}
		vote := &pb.VoteStruct{
			Id:     current["_id"].(primitive.ObjectID).Hex(),
			Video:  current["video"].(primitive.ObjectID).Hex(),
			User:   current["user"].(primitive.ObjectID).Hex(),
			Upvote: current["upvote"].(bool),
		}
		// append document to slice
		votes = append(votes, vote)
	}
	return &pb.ListVotesOfUserResponse{
		Vote: votes,
	}, nil
}
