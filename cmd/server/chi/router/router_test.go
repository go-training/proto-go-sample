package router

import (
	"testing"

	"github.com/go-training/proto-go-sample/pkg/gitea"
	"github.com/go-training/proto-go-sample/pkg/ping"
)

func TestGiteaService(t *testing.T) {
	gitea.MainServiceTest(t, New())
}

func TestPingService(t *testing.T) {
	ping.MainServiceTest(t, New())
}
