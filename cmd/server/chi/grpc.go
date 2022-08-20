package main

import (
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi/v5"
)

func grpcServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
		),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
		),
		compress1KB,
	)

	r.Post(grpcPath+"{name}", grpcHandler(gHandler))
	r.Post(grpcAlphaPath+"{name}", grpcHandler(gAlphaHandler))
}
