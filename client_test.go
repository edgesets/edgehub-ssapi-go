package ssapigo

import (
	"context"
	"fmt"
	sharedpb "github.com/edgesets/edgehub-protocol/shared"
	ssapiv1 "github.com/edgesets/edgehub-protocol/ssapi/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient(t *testing.T) {
	ac, err := NewClient("127.0.0.1:19000", false)
	require.NoError(t, err)
	replies, err := ac.Survey(context.TODO(), &ssapiv1.SurveyRequest{
		Id:            uuid.NewString(),
		ExpectReplies: 1,
		Timeout:       5_000, // millis
		Topic:         "$RPC/http_json",
		Route:         "hello",
		ContentType:   sharedpb.ContentType_JSON.String(),
		PayloadText:   "{\"name\": \"ssapi-go\"}",
	})

	require.NoError(t, err)
	//require.EqualValues(t, replies[0].)
	require.EqualValues(t, 1, len(replies.Results))
	fmt.Println(replies.Results[0].PayloadText)
}
