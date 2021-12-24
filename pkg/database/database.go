package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient interface {
	Connect() error
	Disconnect()
	GetClient() *mongo.Client
}
type mongoClient struct {
	client *mongo.Client
	ctx    *context.Context
}

func CreateNewMongoClient() MongoClient {
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
