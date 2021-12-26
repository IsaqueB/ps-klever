package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VoteModel struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Video  primitive.ObjectID `json:"video" bson:"video"`
	User   primitive.ObjectID `json:"user" bson:"user"`
	Upvote bool               `json:"upvote" bson:"upvote"`
}

type MongoClient interface {
	Connect() error
	Disconnect()
	GetClient() *mongo.Client
}
type mongoClient struct {
	client *mongo.Client
	ctx    *context.Context
}

func NewMongoClient() MongoClient {
	return &mongoClient{}
}

func (mc *mongoClient) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}
	mc.client = client
	mc.ctx = &ctx
	return nil
}

func (mc *mongoClient) Disconnect() {
	mc.client.Disconnect(*mc.ctx)
}

func (mc *mongoClient) GetClient() *mongo.Client {
	return mc.client
}
