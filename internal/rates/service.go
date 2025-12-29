package rates

import (
	"context"
	"errors"
	"fmt"

	"github.com/yosakoo/CurrencyExchanger/internal/currency"
)

type CurrencyService interface {
	GetCurrencyByCode(ctx context.Context, code string) (currency.Currency, error)
}

type Service struct {
	storage         *Storage
	currencyService CurrencyService
}

func NewService(storage *Storage, currencyService CurrencyService) *Service {
	return &Service{
		storage:         storage,
		currencyService: currencyService,
	}
}

func (s *Service) GetExchangeRates(ctx context.Context) ([]ExchangeRate, error) {
	return s.storage.GetExchangeRates(ctx)
}

func (s *Service) GetExchangeRate(ctx context.Context, baseCode, targetCode string) (ExchangeRate, error) {
	return s.storage.GetExchangeRate(ctx, baseCode, targetCode)
}

func (s *Service) CreateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (ExchangeRate, error) {
	if _, err := s.currencyService.GetCurrencyByCode(ctx, baseCode); err != nil {
		return ExchangeRate{}, fmt.Errorf("base currency %s: %w", baseCode, err)
	}

	if _, err := s.currencyService.GetCurrencyByCode(ctx, targetCode); err != nil {
		return ExchangeRate{}, fmt.Errorf("target currency %s: %w", targetCode, err)
	}

	if rate <= 0 {
		return ExchangeRate{}, errors.New("exchange rate must be positive")
	}

	return s.storage.CreateExchangeRate(ctx, baseCode, targetCode, rate)
}

func (s *Service) UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (ExchangeRate, error) {
	if _, err := s.currencyService.GetCurrencyByCode(ctx, baseCode); err != nil {
		return ExchangeRate{}, fmt.Errorf("base currency %s: %w", baseCode, err)
	}

	if _, err := s.currencyService.GetCurrencyByCode(ctx, targetCode); err != nil {
		return ExchangeRate{}, fmt.Errorf("target currency %s: %w", targetCode, err)
	}

	if rate <= 0 {
		return ExchangeRate{}, errors.New("exchange rate must be positive")
	}

	return s.storage.UpdateExchangeRate(ctx, baseCode, targetCode, rate)
}

func (s *Service) CalculateExchange(ctx context.Context, fromCode, toCode string, amount float64) (ExchangeResult, error) {
	if amount <= 0 {
		return ExchangeResult{}, errors.New("amount must be positive")
	}

	fromCurrency, err := s.currencyService.GetCurrencyByCode(ctx, fromCode)
	if err != nil {
		return ExchangeResult{}, fmt.Errorf("source currency %s: %w", fromCode, err)
	}

	toCurrency, err := s.currencyService.GetCurrencyByCode(ctx, toCode)
	if err != nil {
		return ExchangeResult{}, fmt.Errorf("target currency %s: %w", toCode, err)
	}

	if rate, err := s.storage.GetExchangeRate(ctx, fromCode, toCode); err == nil {
		return ExchangeResult{
			BaseCurrency:    rate.BaseCurrency,
			TargetCurrency:  rate.TargetCurrency,
			Rate:            rate.Rate,
			Amount:          amount,
			ConvertedAmount: amount * rate.Rate,
		}, nil
	}

	if rate, err := s.storage.GetExchangeRate(ctx, toCode, fromCode); err == nil {
		reverseRate := 1.0 / rate.Rate

		return ExchangeResult{
			BaseCurrency:    rate.TargetCurrency,
			TargetCurrency:  rate.BaseCurrency,
			Rate:            reverseRate,
			Amount:          amount,
			ConvertedAmount: amount * reverseRate,
		}, nil
	}

	if rate, err := s.calculateCrossRate(ctx, fromCode, toCode); err == nil {
		return ExchangeResult{
			BaseCurrency:    fromCurrency,
			TargetCurrency:  toCurrency,
			Rate:            rate,
			Amount:          amount,
			ConvertedAmount: amount * rate,
		}, nil
	}

	return ExchangeResult{}, fmt.Errorf("no exchange rate available for currency pair %s/%s", fromCode, toCode)
}

func (s *Service) calculateCrossRate(ctx context.Context, fromCode, toCode string) (float64, error) {
	fromToUsdRate, err := s.storage.GetExchangeRate(ctx, fromCode, "USD")
	if err != nil {
		return 0, err
	}

	usdToToRate, err := s.storage.GetExchangeRate(ctx, "USD", toCode)
	if err != nil {
		return 0, err
	}

	return fromToUsdRate.Rate * usdToToRate.Rate, nil
}
