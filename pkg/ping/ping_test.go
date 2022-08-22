package ping

import (
	"net/http"
	"testing"

	"github.com/go-training/proto-go-demo/ping/v1/pingv1connect"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(
		&Service{},
	))
	MainServiceTest(t, mux)
}
