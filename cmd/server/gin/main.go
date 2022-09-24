package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/go-training/proto-go-sample/cmd/server/gin/router"
	"github.com/go-training/proto-go-sample/internal/config"

	otel "github.com/appleboy/go-otel"
	"github.com/appleboy/go-otel/signoz"
	"github.com/appleboy/go-otel/uptrace"
	"github.com/appleboy/graceful"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	var envfile string
	flag.StringVar(&envfile, "env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	_ = godotenv.Load(envfile)
	cfg, err := config.Environ()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("invalid configuration")
	}
	cfg.InitLogging()

	var t otel.TracerProvider

	switch cfg.Otel.ServiceType {
	case "signoz":
		t, err = signoz.New(cfg.Otel.ServiceName, cfg.Otel.CollectorURL, true)
	case "uptrace":
		t = uptrace.New(cfg.Otel.ServiceName)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("can't load otel service")
	}
	defer func() {
		if err := t.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("error shutting down tracer provider")
		}
	}()

	m := graceful.NewManager()

	h := h2c.NewHandler(
		router.New(t.Tracer(), cfg.Otel.ServiceName, cfg.Otel.TargetURL),
		&http2.Server{},
	)
	if cfg.Server.Cert != "" && cfg.Server.Key != "" {
		h = router.New(t.Tracer(), cfg.Otel.ServiceName, cfg.Otel.TargetURL)
	}

	srv := &http.Server{
		Addr:              cfg.Server.Port,
		Handler:           h,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	m.AddRunningJob(func(ctx context.Context) error {
		log.Info().Str("port", cfg.Server.Port).Msg("start http server")
		if err := listenAndServe(srv, cfg.Server.Cert, cfg.Server.Key); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP listen and serve")
		}
		return nil
	})

	m.AddShutdownJob(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("HTTP shutdown")
		}

		return nil
	})

	<-m.Done()
}

func listenAndServe(s *http.Server, certPath string, keyPath string) error {
	if certPath != "" && keyPath != "" {
		return s.ListenAndServeTLS(certPath, keyPath)
	}

	return s.ListenAndServe()
}
