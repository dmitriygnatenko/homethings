package models

import "time"

type ThingNotification struct {
	ThingID          uint64    `db:"thing_id"`
	NotificationDate time.Time `db:"notification_date"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type AddThingNotificationRequest struct {
	ThingID          uint64
	NotificationDate time.Time
}

type UpdateThingNotificationRequest struct {
	ThingID          uint64
	NotificationDate time.Time
}

type ExtThingNotification struct {
	ThingID          uint64    `db:"thing_id"`
	PlaceID          uint64    `db:"place_id"`
	ThingTitle       string    `db:"thing_title"`
	PlaceTitle       string    `db:"place_title"`
	NotificationDate time.Time `db:"notification_date"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
