package gitea

import (
	"net/http"
	"testing"

	"github.com/go-training/proto-go-demo/gitea/v1/giteav1connect"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(giteav1connect.NewGiteaServiceHandler(
		&Service{},
	))
	MainServiceTest(t, mux)
}
