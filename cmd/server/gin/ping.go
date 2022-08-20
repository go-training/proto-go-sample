package main

import (
	"context"
	"fmt"
	"log"

	pingv1 "github.com/go-training/proto-go-demo/ping/v1"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	"github.com/gin-gonic/gin"
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

func pingServiceRoute(r *gin.Engine) {
	compress1KB := connect.WithCompressMinBytes(1024)

	pingService := &PingService{}
	connectPath, connecthandler := pingv1connect.NewPingServiceHandler(
		pingService,
		compress1KB,
	)

	r.POST(connectPath+":name", grpcHandler(connecthandler))
}
