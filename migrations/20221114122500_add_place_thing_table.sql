-- +goose Up
CREATE TABLE place_thing
(
    place_id   INT       NOT NULL,
    thing_id   INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_place_thing_place_id FOREIGN KEY (place_id) REFERENCES place (id),
    CONSTRAINT fk_place_thing_thing_id FOREIGN KEY (thing_id) REFERENCES thing (id)
);

CREATE UNIQUE INDEX idx_unique_place_thing ON place_thing (place_id, thing_id);

CREATE INDEX idx_place_thing_thing ON place_thing (thing_id);

-- +goose Down
DROP TABLE place_thing;
