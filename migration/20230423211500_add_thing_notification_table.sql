-- +goose Up
CREATE TABLE thing_notification
(
    thing_id          INT UNSIGNED NOT NULL PRIMARY KEY,
    notification_date TIMESTAMP NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT `fk_thing_notification_thing_id`
        FOREIGN KEY (thing_id) REFERENCES thing (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- +goose Down
DROP TABLE thing_notification;