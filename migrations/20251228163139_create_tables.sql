-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) NOT NULL,
    fullName VARCHAR(50) NOT NULL,
    sign VARCHAR(10) NOT NULL,
    CONSTRAINT uqCurrenciesCode UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS exchangeRates (
    id SERIAL PRIMARY KEY,
    baseCurrencyId INT NOT NULL,
    targetCurrencyId INT NOT NULL,
    rate DECIMAL(10, 6) NOT NULL,
    CONSTRAINT uqExchangeRatesCurrencyPair UNIQUE (baseCurrencyId, targetCurrencyId),
    CONSTRAINT fkExchangeRatesBaseCurrency FOREIGN KEY (baseCurrencyId) REFERENCES currencies(id) ON DELETE RESTRICT,
    CONSTRAINT fkExchangeRatesTargetCurrency FOREIGN KEY (targetCurrencyId) REFERENCES currencies(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchangeRates;
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd