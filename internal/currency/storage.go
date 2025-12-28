package currency

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type Storage struct {
	conn *sql.DB
}

var (
	ErrCurrencyNotFound = errors.New("currency not found")
	ErrCurrencyExists   = errors.New("currency with this code already exists")
)

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (r *Storage) GetCurrencyByCode(ctx context.Context, code string) (Currency, error) {
	var currency Currency

	err := r.conn.QueryRowContext(ctx, "SELECT id, code, fullname, sign FROM currencies WHERE code = $1", code).Scan(
		&currency.Id,
		&currency.Code,
		&currency.FullName,
		&currency.Sign,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Currency{}, ErrCurrencyNotFound
		}

		return Currency{}, err
	}

	return currency, nil
}

func (r *Storage) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int

	err := r.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM currencies WHERE code = $1", code).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Storage) List(ctx context.Context) ([]Currency, error) {
	rows, err := r.conn.QueryContext(ctx, "SELECT id, code, fullname, sign FROM currencies ORDER BY code")
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error close rows: %v", err)
		}
	}()

	var currencies []Currency

	for rows.Next() {
		var currency Currency

		err := rows.Scan(
			&currency.Id,
			&currency.Code,
			&currency.FullName,
			&currency.Sign,
		)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

func (r *Storage) Create(ctx context.Context, currency *Currency) error {
	err := r.conn.QueryRowContext(ctx,
		`INSERT INTO currencies (code, fullname, sign)
         VALUES ($1, $2, $3)
         ON CONFLICT (code) DO NOTHING
         RETURNING id`,
		currency.Code, currency.FullName, currency.Sign,
	).Scan(&currency.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCurrencyExists
		}

		return err
	}

	return nil
}

func (r *Storage) GetAllCurrencies(ctx context.Context) ([]Currency, error) {
	rows, err := r.conn.QueryContext(ctx, "SELECT id, code, fullname, sign FROM currencies ORDER BY code")
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error close rows: %v", err)
		}
	}()

	var currencies []Currency

	for rows.Next() {
		var currency Currency

		err := rows.Scan(
			&currency.Id,
			&currency.Code,
			&currency.FullName,
			&currency.Sign,
		)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
