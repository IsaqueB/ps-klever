package database_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/IsaqueB/ps-klever/pkg/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateNewMongoClient(t *testing.T) {
	client := database.CreateNewMongoClient()
	switch aux := client.(type) {
	default:
		t.Errorf("Unexpected type. %v", aux)
	case database.MongoClient:
		return
	}
}
func TestConnect(t *testing.T) {
	client := database.CreateNewMongoClient()
	if err := client.Connect(); err != nil {
		t.Errorf("Error while connecting. %v", err)
	}
	defer client.Disconnect()
}
func TestDisconnect(t *testing.T) {
	client := database.CreateNewMongoClient()
	if err := client.Connect(); err != nil {
		t.Errorf("Error while connecting. %v", err)
	}
	client.Disconnect()
	if err := client.GetClient().Ping(context.Background(), nil); err != nil {
		assert.Equal(t, err, mongo.ErrClientDisconnected, "Error different from ErrClientDisconnected")
	} else {
		t.Errorf("Connection still open.")
	}
}
func TestGetClient(t *testing.T) {
	client := database.CreateNewMongoClient()
	if err := client.Connect(); err != nil {
		t.Errorf("Error while connecting. %v", err)
	}
	if err := client.GetClient().Ping(context.Background(), nil); err != nil {
		t.Errorf("Error trying to disconnect with client of GetClient. %v", err)
	}
	defer client.Disconnect()
}

func TestJSONMarshall(t *testing.T) {
	id := primitive.NewObjectID()
	video := primitive.NewObjectID()
	user := primitive.NewObjectID()
	upvote := true
	vote := database.VoteModel{
		ID:     id,
		Video:  video,
		User:   user,
		Upvote: upvote,
	}
	json_str, err := json.Marshal(vote)
	if err != nil {
		t.Errorf("Error formatting the json. %v", err)
	}
	assert.Equal(t, fmt.Sprintf(`{"_id":"%s","video":"%s","user":"%s","upvote":%v}`, id.Hex(), video.Hex(), user.Hex(), upvote), string(json_str))
}

func TestJSONUnmarshall(t *testing.T) {
	id := primitive.NewObjectID()
	video := primitive.NewObjectID()
	user := primitive.NewObjectID()
	upvote := true
	var json_struct = database.VoteModel{}
	json_byte := []byte(fmt.Sprintf(`{"_id":"%s","video":"%s","user":"%s","upvote":%v}`, id.Hex(), video.Hex(), user.Hex(), upvote))
	if err := json.Unmarshal(json_byte, &json_struct); err != nil {
		t.Errorf("Error formatting the struct. %v", err)
	}
	assert.Equal(t, database.VoteModel{ID: id, Video: video, User: user, Upvote: upvote}, json_struct)
}
