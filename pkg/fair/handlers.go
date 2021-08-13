package fair

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPService struct {
	sf StreetFair
}

type errResp struct {
	Msg string `json:"msg"`
}

func prepareResponse(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func errorResponse(w http.ResponseWriter, err error, status int) {
	prepareResponse(w, status)
	_ = json.NewEncoder(w).Encode(&errResp{Msg: err.Error()})
}

func statusByErr(err error) int {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (h *HTTPService) Create(w http.ResponseWriter, r *http.Request) {
	var p Model
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	model, err := h.sf.Create(&p)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	prepareResponse(w, http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&model)
}

func (h *HTTPService) All(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{
		"district":     r.FormValue("district"),
		"region5":      r.FormValue("region5"),
		"name":         r.FormValue("name"),
		"neighborhood": r.FormValue("neighborhood"),
	}
	models, err := h.sf.All(filters)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	prepareResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(&models)
}

func (h *HTTPService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := h.sf.Delete(vars["registry"]); err != nil {
		errorResponse(w, err, statusByErr(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var p Model
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	if vars["registry"] != p.Registry {
		errorResponse(w, errors.New("Cannot update field `registry`"), http.StatusBadRequest)
		return
	}

	if err := h.sf.Update(&p); err != nil {
		errorResponse(w, err, statusByErr(err))
		return
	}
	prepareResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(&p)
}

func (h *HTTPService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	model, err := h.sf.Get(vars["registry"])
	if err != nil {
		errorResponse(w, err, statusByErr(err))
		return
	}
	prepareResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(model)
}

func (h *HTTPService) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/", h.All).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")
	r.HandleFunc("/{registry}/", h.Delete).Methods("DELETE")
	r.HandleFunc("/{registry}/", h.Update).Methods("PUT")
	r.HandleFunc("/{registry}/", h.Get).Methods("GET")
}

func NewHTTPService(sf StreetFair) *HTTPService {
	return &HTTPService{sf: sf}
}
