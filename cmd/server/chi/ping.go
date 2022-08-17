package main

import (
	"context"
	"fmt"
	"log"

	pingv1 "github.com/go-training/proto-go-demo/ping/v1"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi/v5"
)

type PingService struct{}

func (s *PingService) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&pingv1.PingResponse{
		Data: fmt.Sprintf("Hello, %s!", req.Msg.Data),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func pingServiceRoute(r *chi.Mux) {
	compress1KB := connect.WithCompressMinBytes(1024)

	pingService := &PingService{}
	connectPath, connecthandler := pingv1connect.NewPingServiceHandler(
		pingService,
		compress1KB,
	)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(pingv1connect.PingServiceName),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(pingv1connect.PingServiceName),
		compress1KB,
	)

	// grpcHealthCheck
	grpcHealthPath, gHealthHandler := grpchealth.NewHandler(
		grpchealth.NewStaticChecker(pingv1connect.PingServiceName),
		compress1KB,
	)

	r.Post(connectPath+"{name}", grpcHandler(connecthandler))
	r.Post(grpcPath+"{name}", grpcHandler(gHandler))
	r.Post(grpcAlphaPath+"{name}", grpcHandler(gAlphaHandler))
	r.Post(grpcHealthPath+"{name}", grpcHandler(gHealthHandler))
}
