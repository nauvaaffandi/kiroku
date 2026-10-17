-- Migration 006: FX Rates
CREATE TABLE IF NOT EXISTS fx_rates (
    base TEXT NOT NULL,
    quote TEXT NOT NULL,
    rate REAL NOT NULL,
    date DATE NOT NULL,
    source TEXT NOT NULL DEFAULT 'exchangerate.host',
    PRIMARY KEY (base, quote, date)
);

CREATE INDEX IF NOT EXISTS idx_fx_rates_date ON fx_rates(date);