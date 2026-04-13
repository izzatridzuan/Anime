CREATE TABLE users (
    id                   SERIAL PRIMARY KEY,
    email                TEXT NOT NULL UNIQUE,
    name                 TEXT NOT NULL,
    password_hash        TEXT NOT NULL,
    role                 TEXT NOT NULL DEFAULT 'user',
    must_change_password BOOLEAN NOT NULL DEFAULT true,
    created_at           TIMESTAMP NOT NULL DEFAULT NOW()
);