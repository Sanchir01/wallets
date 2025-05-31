-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE operation_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'TRANSFER');

CREATE TABLE IF NOT EXISTS wallets (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    balance DECIMAL(10, 2) DEFAULT 0,
    currency TEXT DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    wallet_id UUID REFERENCES wallets (id) ON DELETE CASCADE,
    sender_wallet_id UUID REFERENCES wallets (id) ON DELETE SET NULL,
    amount DECIMAL(15, 2) NOT NULL,
    type operation_type NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_transfer_logic CHECK (
        (
            type = 'TRANSFER'
            AND sender_wallet_id IS NOT NULL
            AND wallet_id != sender_wallet_id
        )
        OR (
            type IN ('DEPOSIT', 'WITHDRAW')
            AND sender_wallet_id IS NULL
        )
    )
);

CREATE INDEX idx_transactions_wallet_id ON transactions (wallet_id);

CREATE INDEX idx_transactions_sender_wallet_id ON transactions (sender_wallet_id);

CREATE INDEX idx_transactions_type ON transactions (type);

CREATE INDEX idx_transactions_created_at ON transactions (created_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_created_at;

DROP INDEX IF EXISTS idx_transactions_type;

DROP INDEX IF EXISTS idx_transactions_sender_wallet_id;

DROP INDEX IF EXISTS idx_transactions_wallet_id;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS wallets;

-- Drop custom types
DROP TYPE IF EXISTS operation_type;

-- Drop extension (optional, might be used by other tables)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
