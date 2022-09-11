package router

import (
	"log"
	"net/http"
	"time"

	"github.com/go-training/proto-go-sample/internal/grpc"

	"github.com/gin-gonic/gin"
	openapiv2 "github.com/go-training/proto-openapiv2-demo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

func New(serviceName string) *gin.Engine {
	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	r.StaticFS("/public", http.FS(openapiv2.F))

	gRPCRouter(r, grpc.V1Route)
	gRPCRouter(r, grpc.V1AlphaRoute)
	gRPCRouter(r, grpc.HealthRoute)
	gRPCRouter(r, grpc.PingRoute)
	gRPCRouter(r, grpc.GiteaRoute(200*time.Millisecond))

	return r
}
