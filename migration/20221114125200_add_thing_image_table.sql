-- +goose Up
CREATE TABLE thing_image
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    thing_id   INT UNSIGNED NOT NULL,
    image      TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT `fk_thing_image_place_id`
        FOREIGN KEY (thing_id) REFERENCES thing (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- +goose Down
DROP TABLE thing_image;