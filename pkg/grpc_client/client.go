package grpc_client

import (
	"context"
	"net"

	"github.com/IsaqueB/ps-klever/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type Client interface {
	GetClient() proto.VoteClient
	GetConnection() *grpc.ClientConn
	Disconnect() error
	Insert(ctx context.Context, vote *proto.VoteStruct) (string, error)
	Get(ctx context.Context, id string) (*proto.VoteStruct, error)
	UpdateOne(ctx context.Context, id string, new_vote_value bool) (int32, int32, error)
	DeleteOne(ctx context.Context, id string) (int32, error)
}

type client struct {
	vote_c proto.VoteClient
	conn   *grpc.ClientConn
}

func NewGrpcClient(lis *bufconn.Listener) (Client, error) {
	var conn *grpc.ClientConn
	if lis != nil {
		var err error
		conn, err = grpc.DialContext(
			context.Background(),
			"bufnet",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		conn, err = grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
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

func (c *client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Insert(ctx context.Context, vote *proto.VoteStruct) (string, error) {
	response, err := c.vote_c.Insert(ctx, &proto.InsertRequest{Vote: vote})
	if err != nil {
		return "", err
	}
	return response.GetId(), nil
}

func (c *client) Get(ctx context.Context, id string) (*proto.VoteStruct, error) {
	response, err := c.vote_c.Get(ctx, &proto.GetRequest{Id: id})
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
