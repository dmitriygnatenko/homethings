-- +goose Up
CREATE TABLE tag
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title      TEXT      NOT NULL,
    style      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE tag;