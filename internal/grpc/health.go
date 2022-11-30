package grpc

import (
	"net/http"

	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"google.golang.org/grpc/health/grpc_health_v1"
)

func HealthRoute() (string, http.Handler) {
	// grpcHealthCheck
	return grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
			grpc_health_v1.Health_ServiceDesc.ServiceName,
		),
		compress1KB,
	)
}
