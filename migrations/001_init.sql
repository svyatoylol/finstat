CREATE TABLE IF NOT EXISTS transactions (
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT NOT NULL,
    type     TEXT NOT NULL CHECK (type IN ('incoming', 'outgoing')),
    amount   NUMERIC(14, 2) NOT NULL CHECK (amount > 0),
    category TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_id
    ON transactions (user_id, id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_type_id
    ON transactions (user_id, type, id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_category_id
    ON transactions (user_id, category, id);