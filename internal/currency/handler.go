package currency

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// интерфейс на косьюмере при надобности
type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if len(code) != 3 {
		sendError(w, "Currency code must be exactly 3 characters", http.StatusBadRequest)

		return
	}

	upperCode := strings.ToUpper(code)

	currency, err := h.service.GetCurrencyByCode(r.Context(), upperCode)
	if err != nil {
		if errors.Is(err, ErrCurrencyNotFound) {
			sendError(w, "Currency not found", http.StatusNotFound)

			return
		}

		log.Printf("Error get currency: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		log.Printf("Error encoding currency: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.service.List(r.Context())
	if err != nil {
		log.Printf("Error get currencies: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(currencies); err != nil {
		log.Printf("Error encoding currencies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var currency Currency
	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)

		return
	}

	if err := currency.Validate(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)

		return
	}

	currency.Code = strings.ToUpper(currency.Code)

	err := h.service.Create(r.Context(), currency)
	if err != nil {
		log.Printf("Error creating currency: %v", err)

		if errors.Is(err, ErrCurrencyExists) {
			sendError(w, err.Error(), http.StatusConflict)

			return
		}

		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		log.Printf("Error encoding currency: %v", err)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{Message: message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		if _, writeErr := w.Write([]byte(`{"message":"Internal server error"}`)); writeErr != nil {
			log.Printf("Error writing fallback error response: %v", writeErr)
		}
	}
}
