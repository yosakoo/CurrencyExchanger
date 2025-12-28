package currency

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *handler {
	return &handler{storage: storage}
}

func (h *handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if len(code) != 3 {
		http.Error(w, "Currency code must be exactly 3 characters", http.StatusBadRequest)
		return
	}

	upperCode := strings.ToUpper(code)

	currency, err := h.storage.GetCurrencyByCode(r.Context(), upperCode)
	if err != nil {
		if errors.Is(err, ErrCurrencyNotFound) {
			http.Error(w, "Currency not found", http.StatusNotFound)

			return
		}

		log.Printf("Error get currency: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

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
	currencies, err := h.storage.List(r.Context())
	if err != nil {
		log.Printf("Error get currencies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := currency.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currency.Code = strings.ToUpper(currency.Code)

	err := h.storage.Create(r.Context(), &currency)
	if err != nil {
		log.Printf("Error creating currency: %v", err)

		if errors.Is(err, ErrCurrencyExists) {
			http.Error(w, err.Error(), http.StatusConflict)

			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		log.Printf("Error encoding currency: %v", err)
	}
}
