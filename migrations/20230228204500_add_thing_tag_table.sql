-- +goose Up
CREATE TABLE thing_tag
(
    thing_id   INT       NOT NULL,
    tag_id     INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_thing_tag_thing_id FOREIGN KEY (thing_id) REFERENCES thing (id),
    CONSTRAINT fk_thing_tag_tag_id FOREIGN KEY (tag_id) REFERENCES tag (id)
);

CREATE UNIQUE INDEX idx_unique_thing_tag ON thing_tag (thing_id, tag_id);

CREATE INDEX idx_thing_tag_tag ON thing_tag (tag_id);

-- +goose Down
DROP TABLE thing_tag;
