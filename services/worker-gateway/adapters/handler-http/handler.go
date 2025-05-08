package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}
