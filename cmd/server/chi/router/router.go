package router

import (
	"log"
	"net/http"

	"github.com/go-training/proto-go-sample/pkg/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GRPCRouteFn func() (string, http.Handler)

func grpcHandler(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("protocol version:", r.Proto)
		h.ServeHTTP(w, r)
	})
}

func GRPCRouter(r *chi.Mux, fn GRPCRouteFn) {
	p, h := fn()
	r.Post(p+"{name}", grpcHandler(h))
}

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	GRPCRouter(r, grpc.V1Route)
	GRPCRouter(r, grpc.V1AlphaRoute)
	GRPCRouter(r, grpc.PingRoute)
	GRPCRouter(r, grpc.GiteaRoute)
	GRPCRouter(r, grpc.HealthRoute)

	return r
}
