package router

import (
	"context"
	"testing"

	"github.com/go-training/proto-go-sample/internal/gitea"
	"github.com/go-training/proto-go-sample/internal/otel"
	"github.com/go-training/proto-go-sample/internal/ping"
)

func TestGiteaService(t *testing.T) {
	s, err := otel.New("testing", "", true)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Shutdown(context.Background()); err != nil {
			t.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	gitea.MainServiceTest(t, New(s.Tracer(), "testing", ""))
}

func TestPingService(t *testing.T) {
	s, err := otel.New("testing", "", true)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Shutdown(context.Background()); err != nil {
			t.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()
	ping.MainServiceTest(t, New(s.Tracer(), "testing", ""))
}
