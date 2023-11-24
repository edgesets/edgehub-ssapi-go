package ssapigo

import (
	"context"
	sharedpb "github.com/edgesets/edgehub-protocol/shared"
	ssapiv1 "github.com/edgesets/edgehub-protocol/ssapi/v1"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	ssapiv1.ServerSideApiClient
	timeoutMs int32
}

func NewClient(url string, withSSL bool) (*Client, error) {
	var opts []grpc.DialOption
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
		timeoutMs:           5_000,
	}, nil
}

type ReplyDecodeFn func(payload []byte) (proto.Message, error)

func (c *Client) Survey(ctx context.Context, topic string, route string, req proto.Message, waitReplies int) (*ssapiv1.SurveyReply, error) {
	bs, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	return c.ServerSideApiClient.Survey(ctx, &ssapiv1.SurveyRequest{
		Id:           uuid.NewString(),
		WaitReplies:  uint32(waitReplies),
		Timeout:      uint32(c.timeoutMs), // millis
		Topic:        topic,
		Route:        route,
		ContentType:  sharedpb.ContentType_PROTOBUF.String(),
		PayloadBytes: bs,
	})
}
