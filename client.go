package ssapigo

import (
	ssapiv1 "github.com/edgesets/edgehub-protocol/ssapi/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	ssapiv1.ServerSideApiClient
}

func NewClient(url string, withSSL bool) (*Client, error) {
	opts := []grpc.DialOption{}
	if !withSSL {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}
	ac := ssapiv1.NewServerSideApiClient(cc)
	return &Client{
		ServerSideApiClient: ac,
	}, nil
}
