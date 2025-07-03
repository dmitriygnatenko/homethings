-- +goose Up
CREATE TABLE place_image
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    place_id   INT UNSIGNED NOT NULL,
    image      TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT `fk_place_image_place_id`
        FOREIGN KEY (place_id) REFERENCES place (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- +goose Down
DROP TABLE place_image;