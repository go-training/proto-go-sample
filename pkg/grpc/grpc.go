package grpc

import (
	"net/http"

	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var compress1KB = connect.WithCompressMinBytes(1024)

func V1Route() (string, http.Handler) {
	// grpcV1
	return grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
			grpc_health_v1.Health_ServiceDesc.ServiceName,
		),
		compress1KB,
	)
}

func V1AlphaRoute() (string, http.Handler) {
	// grpcV1Alpha
	return grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
			grpc_health_v1.Health_ServiceDesc.ServiceName,
		),
		compress1KB,
	)
}
