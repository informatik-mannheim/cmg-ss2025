package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"gitty.informatik.hs-mannheim.de/steger/cmg-ws202425/services/entity/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {

	h := Handler{service: service, rtr: *mux.NewRouter()}
	h.rtr.HandleFunc("/entity/{id}", h.handleGet).Methods("GET")
	h.rtr.HandleFunc("/entity", h.handleSet).Methods("PUT")
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}

func (h *Handler) handleSet(w http.ResponseWriter, r *http.Request) {
	var entity ports.Entity
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Set(entity, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity, err := h.service.Get(vars["id"], r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}
