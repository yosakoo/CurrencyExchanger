package rates

import (
	"fmt"

	"github.com/yosakoo/CurrencyExchanger/internal/currency"
)

type ExchangeRate struct {
	Id             int               `json:"id"`
	Rate           float64           `json:"rate"`
	BaseCurrency   currency.Currency `json:"baseCurrency"`
	TargetCurrency currency.Currency `json:"targetCurrency"`
}

type ExchangeResult struct {
	BaseCurrency    currency.Currency `json:"baseCurrency"`
	TargetCurrency  currency.Currency `json:"targetCurrency"`
	Rate            float64           `json:"rate"`
	Amount          float64           `json:"amount"`
	ConvertedAmount float64           `json:"convertedAmount"`
}

type CreateExchangeRateRequest struct {
	BaseCurrencyCode   string  `json:"baseCurrencyCode"`
	TargetCurrencyCode string  `json:"targetCurrencyCode"`
	Rate               float64 `json:"rate"`
}

type UpdateExchangeRateRequest struct {
	Rate float64 `json:"rate"`
}

func (r *CreateExchangeRateRequest) Validate() error {
	if r.BaseCurrencyCode == "" {
		return fmt.Errorf("base currency code is required")
	}

	if len(r.BaseCurrencyCode) != 3 {
		return fmt.Errorf("base currency code must be exactly 3 characters")
	}

	if r.TargetCurrencyCode == "" {
		return fmt.Errorf("target currency code is required")
	}

	if len(r.TargetCurrencyCode) != 3 {
		return fmt.Errorf("target currency code must be exactly 3 characters")
	}

	if r.Rate <= 0 {
		return fmt.Errorf("rate must be positive")
	}

	return nil
}

func (r *UpdateExchangeRateRequest) Validate() error {
	if r.Rate <= 0 {
		return fmt.Errorf("rate must be positive")
	}

	return nil
}
