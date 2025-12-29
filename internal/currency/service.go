package currency

import "context"

// делаем сервис для передачи как зависимость в rates
type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetCurrencyByCode(ctx context.Context, code string) (Currency, error) {
	return s.storage.GetCurrencyByCode(ctx, code)
}

func (s *Service) List(ctx context.Context) ([]Currency, error) {
	return s.storage.List(ctx)
}

func (s *Service) Create(ctx context.Context, currency Currency) error {
	return s.storage.Create(ctx, &currency)
}
