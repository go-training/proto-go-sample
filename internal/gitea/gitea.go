package gitea

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
)

type Service struct {
	StreamDelay time.Duration

	giteav1connect.UnimplementedGiteaServiceHandler
}

func (s *Service) Gitea(
	ctx context.Context,
	req *connect.Request[giteav1.GiteaRequest],
) (*connect.Response[giteav1.GiteaResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	log.Println("Te: ", req.Header().Get("Te"))
	log.Println("Grpc-Encoding", req.Header().Get("Grpc-Encoding"))
	log.Println("Grpc-Accept-Encoding", req.Header().Get("Grpc-Accept-Encoding"))
	log.Println("Grpc-Timeout", req.Header().Get("Grpc-Timeout"))
	log.Println("Grpc-Status", req.Header().Get("Grpc-Status"))
	log.Println("Grpc-Message", req.Header().Get("Grpc-Message"))
	log.Println("Grpc-Status-Details-Bin", req.Header().Get("Grpc-Status-Details-Bin"))
	res := connect.NewResponse(&giteav1.GiteaResponse{
		Giteaing: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func (s *Service) Introduce(
	ctx context.Context,
	req *connect.Request[giteav1.IntroduceRequest],
	stream *connect.ServerStream[giteav1.IntroduceResponse],
) error {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	name := req.Msg.Name
	if name == "" {
		name = "Anonymous User"
	}
	intros := []string{name + ", How are you feeling today 01 ?", name + ", How are you feeling today 02 ?"}
	var ticker *time.Ticker
	if s.StreamDelay > 0 {
		ticker = time.NewTicker(s.StreamDelay)
		defer ticker.Stop()
	}
	for _, resp := range intros {
		if ticker != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
		if err := stream.Send(&giteav1.IntroduceResponse{Sentence: resp}); err != nil {
			return err
		}
	}
	return nil
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
	t.Run("gitea", func(t *testing.T) { //nolint: paralleltest
		for _, client := range clients {
			result, err := client.Gitea(context.Background(), connect.NewRequest(&giteav1.GiteaRequest{
				Name: "foobar",
			}))
			assert.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Giteaing)
		}
	})

	t.Run("introduce", func(t *testing.T) { //nolint: paralleltest
		total := 0
		for _, client := range clients {
			request := connect.NewRequest(&giteav1.IntroduceRequest{
				Name: "foobar",
			})
			stream, err := client.Introduce(context.Background(), request)
			assert.Nil(t, err)
			for stream.Receive() {
				total++
			}
			assert.Nil(t, stream.Err())
			assert.Nil(t, stream.Close())
			assert.True(t, total > 0)
		}
	})
}
