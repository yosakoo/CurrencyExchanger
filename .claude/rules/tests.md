---
# Загружается при работе с тестовыми файлами
paths:
  - "**/*_test.go"
---

# Правила для тестов

## Структура
- Unit-тесты рядом с файлом: `payment.go` → `payment_test.go`
- Интеграционные тесты в `internal/repository/integration_test.go`
- Используй table-driven tests для набора похожих случаев

## Моки
- Моки генерируются через mockery: `make mocks`
- Не писать моки вручную
- Моки лежат в `internal/mocks/`

## Что обязательно покрывать
- Happy path
- Граничные значения (0, отрицательные числа, пустые строки)
- Ошибки от зависимостей (БД недоступна, таймаут)
- Бизнес-правила (нельзя списать больше баланса и т.д.)

## Именование
```go
// Паттерн: Test<Функция>_<сценарий>
func TestProcessPayment_Success(t *testing.T) {}
func TestProcessPayment_InsufficientFunds(t *testing.T) {}
func TestProcessPayment_DBError(t *testing.T) {}
```

## Запрещено
- Не использовать `time.Sleep` в тестах — используй моки для времени
- Не обращаться к реальной БД в unit-тестах
- Не игнорировать ошибки через `_`
