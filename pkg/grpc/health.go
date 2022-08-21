package grpc

import (
	"net/http"

	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
)

func HealthRoute() (string, http.Handler) {
	// grpcHealthCheck
	return grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
		),
		compress1KB,
	)
}
