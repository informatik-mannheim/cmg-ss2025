package handler_http

import (
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
)

type Handler struct {
	wg *core.WorkerGatewayService
}

func NewHandler(wg *core.WorkerGatewayService) *Handler {
	return &Handler{wg: wg}
}

// POST /worker/heartbeat
func (h *Handler) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /result
func (h *Handler) SubmitResultHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /register
func (h *Handler) RegisterWorkerHandler(w http.ResponseWriter, r *http.Request) {

}
