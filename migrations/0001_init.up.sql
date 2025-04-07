CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL PRIMARY KEY,
    email        TEXT NOT NULL UNIQUE,
    pass_hash    BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    secret      TEXT NOT NULL UNIQUE
);