-- +goose Up
CREATE TABLE place_image
(
    `id`         INT PRIMARY KEY AUTO_INCREMENT,
    `place_id`   INT          NOT NULL,
    `image`      VARCHAR(255) NOT NULL,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX (`place_id`),
    FOREIGN KEY (`place_id`) REFERENCES place (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE thing_image
(
    `id`         INT PRIMARY KEY AUTO_INCREMENT,
    `thing_id`   INT          NOT NULL,
    `image`      VARCHAR(255) NOT NULL,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX (`thing_id`),
    FOREIGN KEY (`thing_id`) REFERENCES thing (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE place_image;
DROP TABLE thing_image;

