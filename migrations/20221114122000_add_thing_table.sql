-- +goose Up
CREATE TABLE thing
(
    `id`          INT PRIMARY KEY AUTO_INCREMENT,
    `title`       TEXT NOT NULL,
    `description` TEXT,
    `created_at`  DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE thing;

