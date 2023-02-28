-- +goose Up
CREATE TABLE tag
(
    id         SERIAL PRIMARY KEY,
    title      TEXT      NOT NULL,
    style      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE tag;