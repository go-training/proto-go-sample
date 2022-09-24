package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"time"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	pingv1 "github.com/go-training/proto-go-demo/ping/v1"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
	"github.com/go-training/proto-go-sample/internal/config"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func healthCheck(targetURL string, client *http.Client, services ...string) {
	healthClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		targetURL+"grpc.health.v1.Health/Check",
	)

	grpcHealthClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		targetURL+"grpc.health.v1.Health/Check",
		connect.WithGRPC(),
	)

	grpcHealthWebClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		targetURL+"grpc.health.v1.Health/Check",
		connect.WithGRPCWeb(),
	)

	reqClients := []*connect.Client[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse]{
		healthClient,
		grpcHealthClient,
		grpcHealthWebClient,
	}

	for _, n := range services {
		req := &grpc_health_v1.HealthCheckRequest{}
		if n != "" {
			req.Service = n
		}

		for _, c := range reqClients {
			res, err := c.CallUnary(
				context.Background(),
				connect.NewRequest(req),
			)
			if err != nil {
				log.Fatal().Err(err).Msg("err")
			}
			if grpchealth.Status(res.Msg.Status) != grpchealth.StatusServing {
				log.Fatal().Msgf("got status %v, expected %v", res.Msg.Status, grpchealth.StatusServing)
			}
		}
	}
}

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

	c := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	targetURL := "http://" + cfg.Otel.TargetURL + "/"

	connectGiteaClient := giteav1connect.NewGiteaServiceClient(
		c,
		targetURL,
	)

	grpcGiteaClient := giteav1connect.NewGiteaServiceClient(
		c,
		targetURL,
		connect.WithGRPC(),
	)

	grpcWebGiteaClient := giteav1connect.NewGiteaServiceClient(
		c,
		targetURL,
		connect.WithGRPCWeb(),
	)

	giteaClients := []giteav1connect.GiteaServiceClient{
		connectGiteaClient,
		grpcGiteaClient,
		grpcWebGiteaClient,
	}

	connectPingClient := pingv1connect.NewPingServiceClient(
		c,
		targetURL,
	)

	grpcPingClient := pingv1connect.NewPingServiceClient(
		c,
		targetURL,
		connect.WithGRPC(),
	)

	grpcWebPingClient := pingv1connect.NewPingServiceClient(
		c,
		targetURL,
		connect.WithGRPCWeb(),
	)

	pingClients := []pingv1connect.PingServiceClient{
		connectPingClient,
		grpcPingClient,
		grpcWebPingClient,
	}

	for {
		for _, client := range giteaClients {
			req := connect.NewRequest(&giteav1.GiteaRequest{
				Name: "foobar",
			})
			req.Header().Set("Gitea-Header", "hello from connect")
			_, err := client.Gitea(context.Background(), req)
			if err != nil {
				continue
			}
			time.Sleep(100 * time.Millisecond)
		}

		for _, client := range pingClients {
			req := connect.NewRequest(&pingv1.PingRequest{
				Data: "Ping",
			})
			req.Header().Set("Gitea-Header", "hello from connect")
			_, err := client.Ping(context.Background(), req)
			if err != nil {
				continue
			}
		}

		// health check
		healthCheck(
			targetURL,
			c,
			giteav1connect.GiteaServiceName,
			pingv1connect.PingServiceName,
		)
	}
}
