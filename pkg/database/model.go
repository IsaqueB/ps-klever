package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type VoteModel struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Video  primitive.ObjectID `json:"video" bson:"video"`
	User   primitive.ObjectID `json:"user" bson:"user"`
	Upvote bool               `json:"upvote" bson:"upvote"`
}
