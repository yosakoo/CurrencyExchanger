package rates

import "net/http"

func NewRouter(mux *http.ServeMux, service *Service) {
	handler := NewHandler(service)

	mux.HandleFunc("GET /exchangeRates", handler.GetExchangeRates)
	mux.HandleFunc("GET /exchangeRate/{pair}", handler.GetExchangeRate)
	mux.HandleFunc("POST /exchangeRates", handler.CreateExchangeRate)
	mux.HandleFunc("PATCH /exchangeRate/{pair}", handler.UpdateExchangeRate)
	mux.HandleFunc("GET /exchange", handler.Exchange)
}
