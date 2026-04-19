# CurrencyExchanger

REST API для управления валютами и обменными курсами. Написан на Go, хранит данные в PostgreSQL.

## Запуск

```bash
# Поднять БД
docker-compose up -d

# Применить миграции
make migrate-up

# Запустить сервер
go run ./cmd/...
```

## Архитектура

```
cmd/            — точка входа, инициализация зависимостей
internal/
  currency/     — CRUD валют (handler → service → storage)
  rates/        — CRUD курсов + конвертация (handler → service → storage)
  pkg/          — общие утилиты
migrations/     — SQL-миграции через goose
```

Зависимости: `handler → service → storage → PostgreSQL`.

## Стек

- **Go 1.24**, стандартный `net/http` (Go 1.22+ routing)
- **PostgreSQL** через `pgx/v5`
- **goose** — миграции
- **godotenv** — конфиг через `.env`

## API

| Метод | Путь | Описание |
|---|---|---|
| GET | `/currencies` | Список всех валют |
| GET | `/currency/{code}` | Валюта по коду (USD, EUR) |
| POST | `/currency` | Создать валюту |
| GET | `/exchangeRates` | Список всех курсов |
| GET | `/exchangeRate/{pair}` | Курс пары (USDRUB) |
| POST | `/exchangeRates` | Создать курс |
| PATCH | `/exchangeRate/{pair}` | Обновить курс |
| GET | `/exchange?from=USD&to=RUB&amount=100` | Конвертировать |

## Конвертация — алгоритм

`CalculateExchange` пробует три стратегии по порядку:
1. Прямой курс `FROM → TO`
2. Обратный курс `TO → FROM` (инвертирует: `1 / rate`)
3. Кросс-курс через USD: `FROM/USD` × `USD/TO`

Если ни одна не сработала — 404.

## Важные детали

- Коды валют всегда приводятся к верхнему регистру (`usd` → `USD`)
- Кросс-курс работает **только через USD** — другие промежуточные валюты не поддерживаются
- `{pair}` в URL — конкатенация двух 3-буквенных кодов: `USDRUB` = `USD` + `RUB`

## Конфиг (.env)

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=currency_exchanger
DB_USER=postgres
DB_PASSWORD=postgres
```
