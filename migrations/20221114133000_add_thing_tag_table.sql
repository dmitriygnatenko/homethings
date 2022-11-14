-- +goose Up
CREATE TABLE thing_tag
(
    `thing_id` INT NOT NULL,
    `tag_id` INT NOT NULL,
    UNIQUE INDEX (`thing_id`, `tag_id`),
    FOREIGN KEY (`tag_id`) REFERENCES tag (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (`thing_id`) REFERENCES thing (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE thing_tag;
