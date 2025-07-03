-- +goose Up
CREATE TABLE place
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    parent_id  INT,
    title      TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE place;
