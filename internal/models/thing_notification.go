package models

import "time"

type ThingNotification struct {
	ThingID          int
	NotificationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type AddThingNotificationRequest struct {
	ThingID          int
	NotificationDate time.Time
}

type UpdateThingNotificationRequest struct {
	ThingID          int
	NotificationDate time.Time
}

type ExtThingNotification struct {
	ThingID          int
	PlaceID          int
	NotificationDate time.Time
	ThingTitle       string
	PlaceTitle       string
}
