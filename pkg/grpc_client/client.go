package grpc_client

import (
	"context"

	pb "github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
)

type Client interface {
	GetClient() pb.VoteClient
	GetConnection() *grpc.ClientConn
	Disconnect() error
	Insert(ctx context.Context, vote *pb.VoteStruct) (string, error)
	Get(ctx context.Context, id string) (*pb.VoteStruct, error)
	UpdateOne(ctx context.Context, id string, new_vote_value bool) (int32, int32, error)
	DeleteOne(ctx context.Context, id string) (int32, error)
}

type client struct {
	vote_c pb.VoteClient
	conn   *grpc.ClientConn
}

func NewGrpcClient(port string) (Client, error) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	vote_client := pb.NewVoteClient(conn)
	return &client{
		conn:   conn,
		vote_c: vote_client,
	}, nil
}

func (c *client) GetClient() pb.VoteClient {
	return c.vote_c
}

func (c *client) GetConnection() *grpc.ClientConn {
	return c.conn
}

func (c *client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Insert(ctx context.Context, vote *pb.VoteStruct) (string, error) {
	response, err := c.vote_c.Insert(ctx, &pb.InsertRequest{Vote: vote})
	if err != nil {
		return "", err
	}
	return response.GetId(), nil
}

func (c *client) Get(ctx context.Context, id string) (*pb.VoteStruct, error) {
	response, err := c.vote_c.Get(ctx, &pb.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return response.GetVote(), nil
}

func (c *client) UpdateOne(ctx context.Context, id string, new_vote_value bool) (int32, int32, error) {
	response, err := c.vote_c.UpdateOne(ctx, &pb.UpdateOneRequest{Id: id, NewValue: new_vote_value})
	if err != nil {
		return -1, -1, err
	}
	return response.GetMatched(), response.GetModified(), nil
}

func (c *client) DeleteOne(ctx context.Context, id string) (int32, error) {
	response, err := c.vote_c.DeleteOne(ctx, &pb.DeleteOneRequest{Id: id})
	if err != nil {
		return -1, err
	}
	return response.GetDeleted(), nil
}

func (c *client) ListVotesInVideo(ctx context.Context, id string) ([]*pb.VoteStruct, error) {
	response, err := c.vote_c.ListVotesInVideo(ctx, &pb.ListVotesInVideoRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response.GetVote(), nil
}

func (c *client) ListVotesOfUser(ctx context.Context, id string) ([]*pb.VoteStruct, error) {
	response, err := c.vote_c.ListVotesOfUser(ctx, &pb.ListVotesOfUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response.GetVote(), nil
}
