package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
)

func TestGiteaServer(t *testing.T) {
	t.Parallel()
	giteaService := &GiteaServer{}
	mux := http.NewServeMux()
	mux.Handle(giteav1connect.NewGiteaServiceHandler(
		giteaService,
	))
	server := httptest.NewUnstartedServer(mux)
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
