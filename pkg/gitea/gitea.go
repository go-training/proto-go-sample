package gitea

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
)

type Service struct{}

func (s *Service) Gitea(
	ctx context.Context,
	req *connect.Request[giteav1.GiteaRequest],
) (*connect.Response[giteav1.GiteaResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&giteav1.GiteaResponse{
		Giteaing: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func MainServiceTest(t *testing.T, h http.Handler) {
	t.Parallel()
	server := httptest.NewUnstartedServer(h)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := giteav1connect.NewGiteaServiceClient(
		server.Client(),
		server.URL,
	)

	grpcClient := giteav1connect.NewGiteaServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	grpcWebClient := giteav1connect.NewGiteaServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPCWeb(),
	)

	clients := []giteav1connect.GiteaServiceClient{connectClient, grpcClient, grpcWebClient}
	t.Run("gitea", func(t *testing.T) { // nolint: paralleltest
		for _, client := range clients {
			result, err := client.Gitea(context.Background(), connect.NewRequest(&giteav1.GiteaRequest{
				Name: "foobar",
			}))
			assert.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Giteaing)
		}
	})
}
