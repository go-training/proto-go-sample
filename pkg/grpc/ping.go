package grpc

import (
	"net/http"

	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
	"github.com/go-training/proto-go-sample/pkg/ping"
)

func PingRoute() (string, http.Handler) {
	pingService := &ping.Service{}
	return pingv1connect.NewPingServiceHandler(
		pingService,
		compress1KB,
	)
}
