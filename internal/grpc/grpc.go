package grpc

import (
	"net/http"

	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// RouteFn gRPC route registration
type RouteFn func() (string, http.Handler)

var compress1KB = connect.WithCompressMinBytes(1024)

var allServices = []string{
	giteav1connect.GiteaServiceName,
	pingv1connect.PingServiceName,
	grpc_health_v1.Health_ServiceDesc.ServiceName,
}

func V1Route() (string, http.Handler) {
	// grpcV1
	return grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(allServices...),
		compress1KB,
	)
}

func V1AlphaRoute() (string, http.Handler) {
	// grpcV1Alpha
	return grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(allServices...),
		compress1KB,
	)
}
