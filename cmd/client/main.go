package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	giteav1 "github.com/go-training/proto-go-demo/gitea/v1"
	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
	pingv1 "github.com/go-training/proto-go-demo/ping/v1"
	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"golang.org/x/net/http2"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	c := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	connectGiteaClient := giteav1connect.NewGiteaServiceClient(
		c,
		"http://localhost:8080/",
	)

	grpcGiteaClient := giteav1connect.NewGiteaServiceClient(
		c,
		"http://localhost:8080/",
		connect.WithGRPC(),
	)

	grpcHealthClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		c,
		"http://localhost:8080/grpc.health.v1.Health/Check",
		connect.WithGRPC(),
	)

	res, err := grpcHealthClient.CallUnary(
		context.Background(),
		connect.NewRequest(&grpc_health_v1.HealthCheckRequest{}),
	)
	if err != nil {
		log.Fatal(err)
	}
	if grpchealth.Status(res.Msg.Status) != grpchealth.StatusServing {
		log.Fatalf("got status %v, expected %v", res.Msg.Status, grpchealth.StatusServing)
	}

	res, err = grpcHealthClient.CallUnary(
		context.Background(),
		connect.NewRequest(&grpc_health_v1.HealthCheckRequest{
			Service: pingv1connect.PingServiceName,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	if grpchealth.Status(res.Msg.Status) != grpchealth.StatusServing {
		log.Fatalf("got status %v, expected %v", res.Msg.Status, grpchealth.StatusServing)
	}

	giteaClients := []giteav1connect.GiteaServiceClient{connectGiteaClient, grpcGiteaClient}

	for _, client := range giteaClients {
		req := connect.NewRequest(&giteav1.GiteaRequest{
			Name: "foobar",
		})
		req.Header().Set("Gitea-Header", "hello from connect")
		res, err := client.Gitea(context.Background(), req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Message:", res.Msg.Giteaing)
		log.Println("Gitea-Version:", res.Header().Get("Gitea-Version"))
	}

	connectPingClient := pingv1connect.NewPingServiceClient(
		c,
		"http://localhost:8080/",
	)

	grpcPingClient := pingv1connect.NewPingServiceClient(
		c,
		"http://localhost:8080/",
		connect.WithGRPC(),
	)

	pingClients := []pingv1connect.PingServiceClient{connectPingClient, grpcPingClient}

	for _, client := range pingClients {
		req := connect.NewRequest(&pingv1.PingRequest{
			Data: "Ping",
		})
		req.Header().Set("Gitea-Header", "hello from connect")
		res, err := client.Ping(context.Background(), req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Message:", res.Msg.Data)
		log.Println("Gitea-Version:", res.Header().Get("Gitea-Version"))
	}
}
