CREATE TABLE IF NOT EXISTS wallets_incomes (
    income_id UUID PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(wallet_id),
    amount DECIMAL(15,4) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets_expenses (
    income_id UUID PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(wallet_id),
    amount DECIMAL(15,4) NOT NULL,
    created_at TIMESTAMP NOT NULL
);