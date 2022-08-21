package router

import (
	"testing"

	"github.com/go-training/proto-go-sample/pkg/gitea"
)

func TestGiteaService(t *testing.T) {
	gitea.MainServiceTest(t, New())
}
