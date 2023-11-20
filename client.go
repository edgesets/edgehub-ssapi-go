package ssapigo

import (
	"context"
	"fmt"
	sharedpb "github.com/edgesets/edgehub-protocol/shared"
	ssapiv1 "github.com/edgesets/edgehub-protocol/ssapi/v1"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"go.uber.org/multierr"
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

func (c *Client) Survey(ctx context.Context, topic string, route string, req proto.Message, waitReplies int, decodeFn ReplyDecodeFn) ([]proto.Message, error) {
	bs, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	rep, err := c.ServerSideApiClient.Survey(ctx, &ssapiv1.SurveyRequest{
		Id:           uuid.NewString(),
		WaitReplies:  uint32(waitReplies),
		Timeout:      uint32(c.timeoutMs), // millis
		Topic:        topic,
		Route:        route,
		ContentType:  sharedpb.ContentType_PROTOBUF.String(),
		PayloadBytes: bs,
	})
	if e := rep.Error; e != nil {
		return nil, fmt.Errorf("[%d] %s", e.Code, e.Message)
	}
	var errs error
	var replies []proto.Message
	for _, res := range rep.Results {
		if e := res.Error; e != nil {
			errs = multierr.Append(errs, fmt.Errorf("[%d] %s", e.Code, e.Message))
			continue
		}

		// TODO

		//reply, e := decodeFn(res)

	}
	return replies, err
}
