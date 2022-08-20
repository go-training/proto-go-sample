package main

import (
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/gin-gonic/gin"
)

func grpcServiceRoute(r *gin.Engine) {
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

	r.POST(grpcPath+":name", grpcHandler(gHandler))
	r.POST(grpcAlphaPath+":name", grpcHandler(gAlphaHandler))
}
