package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-training/proto-go-sample/cmd/server/gin/router"
	"github.com/go-training/proto-go-sample/internal/otel"

	"github.com/appleboy/graceful"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	listen       int
	certPath     string
	keyPath      string
	serviceName  string
	insecure     bool
	collectorURL string
)

func main() {
	flag.IntVar(&listen, "l", 8080, "listen port")
	flag.StringVar(&certPath, "c", "", "cert path")
	flag.StringVar(&keyPath, "k", "", "key portpath")
	flag.StringVar(&serviceName, "n", "proto-go", "service name")
	flag.BoolVar(&insecure, "insecure", true, "insecure")
	flag.StringVar(&collectorURL, "url", "", "collector URL")
	flag.Parse()

	t, err := otel.New(serviceName, collectorURL, insecure)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := t.TP.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	m := graceful.NewManager()

	h := h2c.NewHandler(
		router.New(serviceName),
		&http2.Server{},
	)
	if certPath != "" && keyPath != "" {
		h = router.New(serviceName)
	}

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(listen),
		Handler:           h,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	m.AddRunningJob(func(ctx context.Context) error {
		log.Println("server listen on port: " + strconv.Itoa(listen))
		if err := listenAndServe(srv, certPath, keyPath); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
		return nil
	})

	m.AddShutdownJob(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP shutdown: %v", err)
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
