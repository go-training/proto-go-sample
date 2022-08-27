package router

import (
	"log"
	"net/http"
	"time"

	"github.com/go-training/proto-go-sample/pkg/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func grpcHandler(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("protocol version:", r.Proto)
		h.ServeHTTP(w, r)
	})
}

func gRPCRouter(r *chi.Mux, fn grpc.RouteFn) {
	p, h := fn()
	r.Post(p+"{name}", grpcHandler(h))
}

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})
	gRPCRouter(r, grpc.V1Route)
	gRPCRouter(r, grpc.V1AlphaRoute)
	gRPCRouter(r, grpc.HealthRoute)
	gRPCRouter(r, grpc.PingRoute)
	gRPCRouter(r, grpc.GiteaRoute(200*time.Millisecond))

	return r
}
