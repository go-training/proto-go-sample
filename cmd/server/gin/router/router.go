package router

import (
	"log"
	"net/http"
	"time"

	"github.com/go-training/proto-go-sample/internal/grpc"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	openapiv2 "github.com/go-training/proto-openapiv2-demo"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
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

func New(t trace.Tracer, serviceName, targetURL string) *gin.Engine {
	r := gin.New()
	r.Use(otelgin.Middleware(serviceName))
	r.Use(requestid.New())
	r.Use(logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.With().
				Str("request_id", requestid.Get(c)).
				Str("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()).
				Logger()
		})))
	r.StaticFS("/public", http.FS(openapiv2.F))

	gRPCRouter(r, grpc.V1Route)
	gRPCRouter(r, grpc.V1AlphaRoute)
	gRPCRouter(r, grpc.HealthRoute)
	gRPCRouter(r, grpc.PingRoute)
	gRPCRouter(r, grpc.GiteaRoute(t, targetURL, 200*time.Millisecond))

	return r
}
