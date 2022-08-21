package grpc

import (
	"net/http"

	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-sample/pkg/gitea"
)

func GiteaRoute() (string, http.Handler) {
	giteaService := &gitea.Service{}

	return giteav1connect.NewGiteaServiceHandler(
		giteaService,
		compress1KB,
	)
}
