CREATE TABLE IF NOT EXISTS movements (
    id              UUID PRIMARY KEY,
    account_id      VARCHAR(255) NOT NULL,
    institution_id  VARCHAR(255) NOT NULL,
    description     TEXT,
    amount          DECIMAL(10, 2) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    date            TIMESTAMP WITH TIME ZONE NOT NULL,
    source          VARCHAR(50) NOT NULL,
    category        VARCHAR(100),
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_movements_account_id ON movements (account_id);