CREATE TABLE IF NOT EXISTS urls (
    id bigserial PRIMARY KEY,
    original text NOT NULL UNIQUE,
    alias text NOT NULL UNIQUE
);
