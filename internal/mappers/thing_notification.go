package mappers

import (
	"time"

	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func ToAddThingNotificationRequest(req dto.AddThingNotificationRequest) (*models.AddThingNotificationRequest, error) {
	date, err := time.Parse(time.RFC3339, req.NotificationDate)
	if err != nil {
		return nil, err
	}

	return &models.AddThingNotificationRequest{
		ThingID:          req.ThingID,
		NotificationDate: date,
	}, nil
}

func ToUpdateThingNotificationRequest(
	thingID uint64,
	req dto.UpdateThingNotificationRequest,
) (*models.UpdateThingNotificationRequest, error) {
	date, err := time.Parse(time.RFC3339, req.NotificationDate)
	if err != nil {
		return nil, err
	}

	return &models.UpdateThingNotificationRequest{
		ThingID:          thingID,
		NotificationDate: date,
	}, nil
}

func ToThingNotificationResponse(req models.ThingNotification) dto.ThingNotificationResponse {
	return dto.ThingNotificationResponse{
		ThingID:          req.ThingID,
		NotificationDate: req.NotificationDate.Format(defaultDateTimeLayout),
		CreatedAt:        req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt:        req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ToThingNotificationExtResponse(req models.ExtThingNotification) dto.ThingNotificationExtResponse {
	return dto.ThingNotificationExtResponse{
		ThingID:          req.ThingID,
		PlaceID:          req.PlaceID,
		ThingTitle:       req.ThingTitle,
		PlaceTitle:       req.PlaceTitle,
		NotificationDate: req.NotificationDate.Format(defaultDateTimeLayout),
		CreatedAt:        req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt:        req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ToThingNotificationsExtResponse(req []models.ExtThingNotification) dto.ThingNotificationsExtResponse {
	res := make([]dto.ThingNotificationExtResponse, 0, len(req))

	for _, notification := range req {
		res = append(res, ToThingNotificationExtResponse(notification))
	}

	return dto.ThingNotificationsExtResponse{Notifications: res}
}
