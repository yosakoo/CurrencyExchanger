package rates

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/yosakoo/CurrencyExchanger/internal/currency"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	exchangeRates, err := h.service.GetExchangeRates(r.Context())
	if err != nil {
		log.Printf("Error get exchange rates: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(exchangeRates); err != nil {
		log.Printf("Error encoding exchange rates: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("pair")

	baseCode, targetCode, err := validateCurrencyCodes(pathValue)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	exchangeRate, err := h.service.GetExchangeRate(r.Context(), baseCode, targetCode)
	if err != nil {
		if errors.Is(err, ErrExchangeRateNotFound) {
			sendError(w, "Exchange rate not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, currency.ErrCurrencyNotFound) {
			sendError(w, "Currency not found", http.StatusNotFound)
			return
		}

		log.Printf("Error get exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(exchangeRate); err != nil {
		log.Printf("Error encoding exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) CreateExchangeRate(w http.ResponseWriter, r *http.Request) {
	var req CreateExchangeRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	exchangeRate, err := h.service.CreateExchangeRate(r.Context(), req.BaseCurrencyCode, req.TargetCurrencyCode, req.Rate)
	if err != nil {
		if errors.Is(err, currency.ErrCurrencyNotFound) {
			sendError(w, "One or both currencies not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, ErrExchangeRateExists) {
			sendError(w, err.Error(), http.StatusConflict)
			return
		}

		log.Printf("Error creating exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(exchangeRate); err != nil {
		log.Printf("Error encoding exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}
}

func (h *handler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("pair")

	baseCode, targetCode, err := validateCurrencyCodes(pathValue)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req UpdateExchangeRateRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	exchangeRate, err := h.service.UpdateExchangeRate(r.Context(), baseCode, targetCode, req.Rate)
	if err != nil {
		if errors.Is(err, ErrExchangeRateNotFound) {
			sendError(w, "Exchange rate not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, currency.ErrCurrencyNotFound) {
			sendError(w, "Currency not found", http.StatusNotFound)
			return
		}

		log.Printf("Error updating exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(exchangeRate); err != nil {
		log.Printf("Error encoding exchange rate: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) Exchange(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if from == "" || to == "" || amountStr == "" {
		sendError(w, "Missing required query parameters: from, to, amount", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		sendError(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	if amount <= 0 {
		sendError(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	result, err := h.service.CalculateExchange(r.Context(), from, to, amount)
	if err != nil {
		if errors.Is(err, ErrExchangeRateNotFound) || errors.Is(err, currency.ErrCurrencyNotFound) {
			sendError(w, "Exchange rate not available for this currency pair", http.StatusNotFound)
			return
		}

		log.Printf("Error calculating exchange: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding exchange result: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)

		return
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{Message: message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		if _, writeErr := w.Write([]byte(`{"message":"Internal server error"}`)); writeErr != nil {
			log.Printf("Error writing fallback error response: %v", writeErr)
		}
	}
}

func validateCurrencyCodes(pair string) (string, string, error) {
	if len(pair) != 6 {
		return "", "", fmt.Errorf("currency pair codes must be exactly 6 characters")
	}

	baseCode := strings.ToUpper(pair[:3])
	targetCode := strings.ToUpper(pair[3:])

	for _, code := range []string{baseCode, targetCode} {
		for _, char := range code {
			if char < 'A' || char > 'Z' {
				return "", "", fmt.Errorf("currency codes must contain only English letters (A-Z)")
			}
		}
	}

	return baseCode, targetCode, nil
}
