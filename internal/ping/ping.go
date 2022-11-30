package ping

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pingv1 "github.com/go-training/proto-go-demo/ping/v1"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
)

type Service struct {
	pingv1connect.UnimplementedPingServiceHandler
}

func (s *Service) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	res := connect.NewResponse(&pingv1.PingResponse{
		Data: "pong",
	})
	return res, nil
}

func MainServiceTest(t *testing.T, h http.Handler) {
	t.Parallel()
	server := httptest.NewUnstartedServer(h)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
	)

	grpcClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	grpcWebClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPCWeb(),
	)

	clients := []pingv1connect.PingServiceClient{connectClient, grpcClient, grpcWebClient}
	t.Run("gitea", func(t *testing.T) { //nolint: paralleltest
		for _, client := range clients {
			result, err := client.Ping(context.Background(), connect.NewRequest(&pingv1.PingRequest{
				Data: "foobar",
			}))
			assert.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Data)
		}
	})
}
