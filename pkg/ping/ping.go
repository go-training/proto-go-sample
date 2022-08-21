package ping

import (
	"context"
	"fmt"
	"log"

	pingv1 "github.com/go-training/proto-go-demo/ping/v1"

	"github.com/bufbuild/connect-go"
)

type Service struct{}

func (s *Service) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&pingv1.PingResponse{
		Data: fmt.Sprintf("Hello, %s!", req.Msg.Data),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}
