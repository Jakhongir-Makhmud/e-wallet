CREATE TABLE IF NOT EXISTS wallets_income (
    income_id UUID PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(wallet_id),
    created_at TIMESTAMP NOT NULL,
);