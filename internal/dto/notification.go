package dto

type AddThingNotificationRequest struct {
	ThingID          int    `json:"thing_id" validate:"required"`
	NotificationDate string `json:"notification_date" validate:"required" format:"date-time"`
}

type UpdateThingNotificationRequest struct {
	NotificationDate string `json:"notification_date" validate:"required" format:"date-time"`
}

type ThingNotificationResponse struct {
	ThingID          int    `json:"id"`
	NotificationDate string `json:"notification_date"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type ThingNotificationsExtResponse struct {
	Notifications []ThingNotificationExtResponse `json:"notifications"`
}

type ThingNotificationExtResponse struct {
	ThingID          int    `json:"thing_id"`
	PlaceID          int    `json:"place_id"`
	ThingTitle       string `json:"thing_title"`
	PlaceTitle       string `json:"place_title"`
	NotificationDate string `json:"notification_date"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
