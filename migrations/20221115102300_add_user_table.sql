-- +goose Up
CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    username   TEXT      NOT NULL UNIQUE,
    password   TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE "user";