package currency

import "net/http"

func NewRouter(mux *http.ServeMux, service *Service) {
	handler := NewHandler(service)

	mux.HandleFunc("GET /currency/{code}", handler.GetByCode)
	mux.HandleFunc("GET /currencies", handler.List)
	mux.HandleFunc("POST /currencies", handler.Create)
}
