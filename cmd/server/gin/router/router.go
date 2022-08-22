package router

import (
	"log"
	"net/http"

	"github.com/go-training/proto-go-sample/pkg/grpc"

	"github.com/gin-gonic/gin"
)

func grpcHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("protocol version:", c.Request.Proto)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func gRPCRouter(r *gin.Engine, fn grpc.RouteFn) {
	p, h := fn()
	r.POST(p+":name", grpcHandler(h))
}

func New() *gin.Engine {
	r := gin.Default()
	gRPCRouter(r, grpc.V1Route)
	gRPCRouter(r, grpc.V1AlphaRoute)
	gRPCRouter(r, grpc.HealthRoute)
	gRPCRouter(r, grpc.PingRoute)
	gRPCRouter(r, grpc.GiteaRoute)

	return r
}
