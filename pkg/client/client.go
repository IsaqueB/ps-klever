package client

import (
	"context"

	"github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
)

type Client interface {
	GetClient() proto.VoteClient
	GetConnection() *grpc.ClientConn
	Insert(ctx context.Context, vote *proto.VoteStruct) (string, error)
	Get(ctx context.Context, id string) (*proto.VoteStruct, error)
	UpdateOne(ctx context.Context, id string, new_vote_value bool) (int32, int32, error)
	DeleteOne(ctx context.Context, id string) (int32, error)
}

type client struct {
	vote_c proto.VoteClient
	conn   *grpc.ClientConn
}

func NewGrpcClient() (Client, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	vote_client := proto.NewVoteClient(conn)
	return &client{
		conn:   conn,
		vote_c: vote_client,
	}, nil
}

func (c *client) GetClient() proto.VoteClient {
	return c.vote_c
}

func (c *client) GetConnection() *grpc.ClientConn {
	return c.conn
}

func (c *client) Insert(ctx context.Context, vote *proto.VoteStruct) (string, error) {
	response, err := c.vote_c.Insert(ctx, &proto.InsertRequest{Vote: vote})
	if err != nil {
		return "", err
	}
	return response.GetId(), nil
}

func (c *client) Get(ctx context.Context, id string) (*proto.VoteStruct, error) {
	response, err := c.vote_c.Get(ctx, &proto.GetRequest{Id: "61c4014dd6f4074658db9772"})
	if err != nil {
		return nil, err
	}
	return response.GetVote(), nil
}

func (c *client) UpdateOne(ctx context.Context, id string, new_vote_value bool) (int32, int32, error) {
	response, err := c.vote_c.UpdateOne(ctx, &proto.UpdateOneRequest{Id: id, NewValue: new_vote_value})
	if err != nil {
		return -1, -1, err
	}
	return response.GetMatched(), response.GetModified(), nil
}

func (c *client) DeleteOne(ctx context.Context, id string) (int32, error) {
	response, err := c.vote_c.DeleteOne(ctx, &proto.DeleteOneRequest{Id: id})
	if err != nil {
		return -1, err
	}
	return response.GetDeleted(), nil
}

// 	// response3, err := c.GetAllVotesOfVideo(context.Background(), &proto.Message{Body: "61c29a8131e2495f3a6c29a4"})
// 	// if err != nil {
// 	// 	log.Fatalf("could not connect! %v", err)
// 	// }
// 	// log.Printf("Response from Server: %s", response3)

// 	stream, err := c.ListVotesInVideo(context.Background(), &proto.ListVotesInVideoRequest{Id: "61c29a8131e2495f3a6c29a4"})
// 	if err != nil {
// 		log.Fatalf("%v", err)
// 	}
// 	for {
// 		res, err := stream.Recv()
// 		if err == io.EOF {
// 			return
// 		}
// 		if err != nil {
// 			log.Fatalf("%v", err)
// 		}
// 		recv_vote := res.GetVote()
// 		fmt.Printf("Vote: %v\n", recv_vote)
// 	}
// }
