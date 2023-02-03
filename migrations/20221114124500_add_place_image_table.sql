-- +goose Up
CREATE TABLE place_image
(
    id         SERIAL PRIMARY KEY,
    place_id   INT       NOT NULL,
    image      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_place_image_place_id FOREIGN KEY (place_id) REFERENCES place (id)
);

CREATE INDEX idx_place_image_place ON place_image (place_id);

-- +goose Down
DROP TABLE place_image;