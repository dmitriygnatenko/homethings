-- +goose Up
CREATE TABLE place_thing
(
    place_id   INT UNSIGNED NOT NULL,
    thing_id   INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (place_id, thing_id),
    CONSTRAINT `fk_place_thing_place_id`
        FOREIGN KEY (place_id) REFERENCES place (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_place_thing_thing_id`
        FOREIGN KEY (thing_id) REFERENCES thing (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- +goose Down
DROP TABLE place_thing;
