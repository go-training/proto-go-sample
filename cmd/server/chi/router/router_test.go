package router

import (
	"testing"

	"github.com/go-training/proto-go-sample/internal/gitea"
	"github.com/go-training/proto-go-sample/internal/ping"
)

func TestGiteaService(t *testing.T) {
	gitea.MainServiceTest(t, New("testing"))
}

func TestPingService(t *testing.T) {
	ping.MainServiceTest(t, New("testing"))
}
