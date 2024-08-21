package models

import "time"

type ThingNotification struct {
	ThingID          uint64
	NotificationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
	ThingID          uint64
	PlaceID          uint64
	ThingTitle       string
	PlaceTitle       string
	NotificationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
