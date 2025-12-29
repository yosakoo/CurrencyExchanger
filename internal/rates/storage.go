package rates

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/yosakoo/CurrencyExchanger/internal/currency"
)

var (
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
	ErrExchangeRateExists   = errors.New("exchange rate for this currency pair already exists")
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (r *Storage) GetExchangeRates(ctx context.Context) ([]ExchangeRate, error) {
	query := `
		SELECT er.id, er.rate,
			   bc.id, bc.code, bc.fullname, bc.sign,
			   tc.id, tc.code, tc.fullname, tc.sign
		FROM exchangerates er
		JOIN currencies bc ON er.basecurrencyid = bc.id
		JOIN currencies tc ON er.targetcurrencyid = tc.id
		ORDER BY er.id
	`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error close rows: %v", closeErr)
		}
	}()

	var exchangeRates []ExchangeRate

	for rows.Next() {
		var er ExchangeRate

		err := rows.Scan(
			&er.Id, &er.Rate,
			&er.BaseCurrency.Id, &er.BaseCurrency.Code, &er.BaseCurrency.FullName, &er.BaseCurrency.Sign,
			&er.TargetCurrency.Id, &er.TargetCurrency.Code, &er.TargetCurrency.FullName, &er.TargetCurrency.Sign,
		)
		if err != nil {
			return nil, err
		}

		exchangeRates = append(exchangeRates, er)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exchangeRates, nil
}

func (r *Storage) GetExchangeRate(ctx context.Context, baseCode, targetCode string) (ExchangeRate, error) {
	query := `
		SELECT er.id, er.rate,
			   bc.id, bc.code, bc.fullname, bc.sign,
			   tc.id, tc.code, tc.fullname, tc.sign
		FROM exchangerates er
		JOIN currencies bc ON er.basecurrencyid = bc.id
		JOIN currencies tc ON er.targetcurrencyid = tc.id
		WHERE bc.code = $1 AND tc.code = $2
	`

	var er ExchangeRate

	err := r.conn.QueryRowContext(ctx, query, baseCode, targetCode).Scan(
		&er.Id, &er.Rate,
		&er.BaseCurrency.Id, &er.BaseCurrency.Code, &er.BaseCurrency.FullName, &er.BaseCurrency.Sign,
		&er.TargetCurrency.Id, &er.TargetCurrency.Code, &er.TargetCurrency.FullName, &er.TargetCurrency.Sign,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ExchangeRate{}, ErrExchangeRateNotFound
		}
		return ExchangeRate{}, err
	}

	return er, nil
}

func (r *Storage) CreateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (ExchangeRate, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return ExchangeRate{}, err
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("Error rolling back transaction: %v", rbErr)
		}
	}()

	var baseId, targetId int
	if scanErr := tx.QueryRowContext(ctx, "SELECT id FROM currencies WHERE code = $1", baseCode).Scan(&baseId); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return ExchangeRate{}, currency.ErrCurrencyNotFound
		}
		return ExchangeRate{}, scanErr
	}

	if scanErr := tx.QueryRowContext(ctx, "SELECT id FROM currencies WHERE code = $1", targetCode).Scan(&targetId); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return ExchangeRate{}, currency.ErrCurrencyNotFound
		}
		return ExchangeRate{}, scanErr
	}

	var erId int

	err = tx.QueryRowContext(ctx,
		`INSERT INTO exchangerates (basecurrencyid, targetcurrencyid, rate) 
		 VALUES ($1, $2, $3) 
		 ON CONFLICT (basecurrencyid, targetcurrencyid) DO NOTHING
		 RETURNING id`,
		baseId, targetId, rate,
	).Scan(&erId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ExchangeRate{}, ErrExchangeRateExists
		}
		return ExchangeRate{}, err
	}

	if err = tx.Commit(); err != nil {
		return ExchangeRate{}, err
	}

	return r.GetExchangeRate(ctx, baseCode, targetCode)
}

func (r *Storage) UpdateExchangeRate(ctx context.Context, baseCode, targetCode string, rate float64) (ExchangeRate, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return ExchangeRate{}, err
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("Error rolling back transaction: %v", rbErr)
		}
	}()

	var baseId, targetId int
	if scanErr := tx.QueryRowContext(ctx, "SELECT id FROM currencies WHERE code = $1", baseCode).Scan(&baseId); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return ExchangeRate{}, currency.ErrCurrencyNotFound
		}
		return ExchangeRate{}, scanErr
	}

	if scanErr := tx.QueryRowContext(ctx, "SELECT id FROM currencies WHERE code = $1", targetCode).Scan(&targetId); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return ExchangeRate{}, currency.ErrCurrencyNotFound
		}
		return ExchangeRate{}, scanErr
	}

	result, err := tx.ExecContext(ctx,
		"UPDATE exchangerates SET rate = $1 WHERE basecurrencyid = $2 AND targetcurrencyid = $3",
		rate, baseId, targetId,
	)
	if err != nil {
		return ExchangeRate{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ExchangeRate{}, err
	}

	if rowsAffected == 0 {
		return ExchangeRate{}, ErrExchangeRateNotFound
	}

	if err = tx.Commit(); err != nil {
		return ExchangeRate{}, err
	}

	return r.GetExchangeRate(ctx, baseCode, targetCode)
}
