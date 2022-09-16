package grpc

import (
	"net/http"
	"time"

	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	"github.com/go-training/proto-go-sample/internal/gitea"

	"go.opentelemetry.io/otel/trace"
)

func GiteaRoute(t trace.Tracer, d time.Duration) RouteFn {
	giteaService := &gitea.Service{
		StreamDelay: d,
		Tracer:      t,
	}

	return func() (string, http.Handler) {
		return giteav1connect.NewGiteaServiceHandler(
			giteaService,
			compress1KB,
		)
	}
}
