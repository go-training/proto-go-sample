package main

import (
	"context"
	"fmt"
	"log"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi/v5"
)

type GiteaServer struct{}

func (s *GiteaServer) Gitea(
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

func giteaServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	giteaService := &GiteaServer{}
	connectPath, connecthandler := giteav1connect.NewGiteaServiceHandler(
		giteaService,
		compress1KB,
	)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(giteav1connect.GiteaServiceName),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(giteav1connect.GiteaServiceName),
		compress1KB,
	)

	r.Post(connectPath+"{name}", grpcHandler(connecthandler))
	r.Post(grpcPath+"{name}", grpcHandler(gHandler))
	r.Post(grpcAlphaPath+"{name}", grpcHandler(gAlphaHandler))
}
