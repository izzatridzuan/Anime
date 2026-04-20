CREATE TABLE service_accounts(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    api_key_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);