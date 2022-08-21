package router

import (
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-sample/pkg/gitea"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
)

func giteaServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	giteaService := &gitea.Service{}
	connectPath, connecthandler := giteav1connect.NewGiteaServiceHandler(
		giteaService,
		compress1KB,
	)

	r.Post(connectPath+"{name}", grpcHandler(connecthandler))
}
