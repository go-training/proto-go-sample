package router

import (
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
	"github.com/go-training/proto-go-sample/pkg/ping"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
)

func pingServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	pingService := &ping.Service{}
	connectPath, connecthandler := pingv1connect.NewPingServiceHandler(
		pingService,
		compress1KB,
	)

	r.Post(connectPath+"{name}", grpcHandler(connecthandler))
}
