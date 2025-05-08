package core

import (
	//"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type WorkerGatewayService struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewWorkerGatewayService(repo ports.Repo, notifier ports.Notifier) *WorkerGatewayService {
	return &WorkerGatewayService{
		repo:     repo,
		notifier: notifier,
	}
}

var _ ports.Api = (*WorkerGatewayService)(nil)
