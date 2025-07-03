-- +goose Up
CREATE TABLE thing_tag
(
    thing_id   INT UNSIGNED NOT NULL,
    tag_id     INT UNSIGNED NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    UNIQUE (thing_id, tag_id),
    CONSTRAINT `fk_thing_tag_thing_id`
        FOREIGN KEY (thing_id) REFERENCES thing (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_thing_tag_tag_id`
        FOREIGN KEY (tag_id) REFERENCES tag (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- +goose Down
DROP TABLE thing_tag;
