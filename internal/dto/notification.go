package dto

type AddThingNotificationRequest struct {
	ThingID          int    `json:"thing_id" validate:"required"`
	NotificationDate string `json:"notification_date" validate:"required"`
}

type UpdateThingNotificationRequest struct {
	ThingID          int    `json:"thing_id" validate:"required"`
	NotificationDate string `json:"notification_date" validate:"required"`
}
