package handler_http

/*

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {

	h := Handler{service: service, rtr: *mux.NewRouter()}
	h.rtr.HandleFunc("/job/{id}", h.handleGet).Methods("GET")
	h.rtr.HandleFunc("/job", h.handleSet).Methods("PUT")
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}

func (h *Handler) handleSet(w http.ResponseWriter, r *http.Request) {
	var job ports.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Set(job, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job, err := h.service.Get(vars["id"], r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}


*/
