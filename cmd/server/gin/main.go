package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-training/proto-go-sample/cmd/server/gin/router"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	listen   int
	certPath string
	keyPath  string
)

func main() {
	flag.IntVar(&listen, "l", 8080, "listen port")
	flag.StringVar(&certPath, "c", "", "cert path")
	flag.StringVar(&keyPath, "k", "", "key portpath")
	flag.Parse()

	h := h2c.NewHandler(
		router.New(),
		&http2.Server{},
	)
	if certPath != "" && keyPath != "" {
		h = router.New()
	}

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(listen),
		Handler:           h,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := listenAndServe(srv, certPath, keyPath); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()
	log.Println("server listen on port: " + strconv.Itoa(listen))

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err) // nolint:gocritic
	}
}

func listenAndServe(s *http.Server, certPath string, keyPath string) error {
	if certPath != "" && keyPath != "" {
		return s.ListenAndServeTLS(certPath, keyPath)
	}

	return s.ListenAndServe()
}
