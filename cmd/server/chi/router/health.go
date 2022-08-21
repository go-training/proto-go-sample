package router

import (
	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
)

func healthServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	// grpcHealthCheck
	grpcHealthPath, gHealthHandler := grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
		),
		compress1KB,
	)

	r.Post(grpcHealthPath+"{name}", grpcHandler(gHealthHandler))
}
