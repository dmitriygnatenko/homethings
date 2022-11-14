-- +goose Up
CREATE TABLE place
(
    `id`         INT PRIMARY KEY AUTO_INCREMENT,
    `parent_id`  INT,
    `title`      TEXT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE place;
