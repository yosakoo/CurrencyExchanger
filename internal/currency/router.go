package currency

import "net/http"

func NewRouter(storage *Storage) http.Handler {
	mux := http.NewServeMux()
	handler := NewHandler(storage)

	mux.HandleFunc("GET /currency/{code}", handler.GetByCode)
	mux.HandleFunc("GET /currencies", handler.List)
	mux.HandleFunc("POST /currencies", handler.Create)

	return mux
}
