package client

import (
	"context"

	pb "github.com/reshmavatkar/kv-store/generated" // Protobuf-generated Go code

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------- gRPC Client Interface -----------

type StoreClient interface {
	Put(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type grpcStoreClient struct {
	conn   *grpc.ClientConn
	client pb.KeyValueStoreClient
}

func NewStoreClient(addr string) (StoreClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewKeyValueStoreClient(conn)
	return &grpcStoreClient{conn: conn, client: client}, nil
}

func (c *grpcStoreClient) Put(ctx context.Context, key, value string) error {
	_, err := c.client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
	return err
}

func (c *grpcStoreClient) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func (c *grpcStoreClient) Delete(ctx context.Context, key string) error {
	_, err := c.client.Delete(ctx, &pb.DeleteRequest{Key: key})
	return err
}

func (c *grpcStoreClient) Close() error {
	return c.conn.Close()
}
