-- +goose Up
CREATE TABLE thing_image
(
    id         SERIAL PRIMARY KEY,
    thing_id   INT       NOT NULL,
    image      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_thing_image_thing_id FOREIGN KEY (thing_id) REFERENCES thing (id)
);

CREATE INDEX idx_thing_image_thing ON thing_image (thing_id);

-- +goose Down
DROP TABLE thing_image;