-- +goose Up
CREATE TABLE place_thing
(
    `place_id` INT NOT NULL,
    `thing_id` INT NOT NULL,
    UNIQUE INDEX (`place_id`, `thing_id`),
    FOREIGN KEY (`place_id`) REFERENCES place (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (`thing_id`) REFERENCES thing (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE place_thing;
