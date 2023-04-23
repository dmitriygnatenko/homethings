package mappers

import (
	"time"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToAddThingNotificationRequestModel(req dto.AddThingNotificationRequest) (*models.AddThingNotificationRequest, error) {
	date, err := time.Parse(time.RFC3339, req.NotificationDate)
	if err != nil {
		return nil, err
	}

	return &models.AddThingNotificationRequest{
		ThingID:          req.ThingID,
		NotificationDate: date,
	}, nil
}

func ConvertToUpdateThingNotificationRequestModel(thingID int, req dto.UpdateThingNotificationRequest) (*models.UpdateThingNotificationRequest, error) {
	date, err := time.Parse(time.RFC3339, req.NotificationDate)
	if err != nil {
		return nil, err
	}

	return &models.UpdateThingNotificationRequest{
		ThingID:          thingID,
		NotificationDate: date,
	}, nil
}

func ConvertToThingNotificationResponseDTO(req models.ThingNotification) dto.ThingNotificationResponse {
	return dto.ThingNotificationResponse{
		ThingID:          req.ThingID,
		NotificationDate: req.NotificationDate.Format(defaultDateTimeLayout),
		CreatedAt:        req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt:        req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ConvertToThingNotificationExtResponseDTO(req models.ExtThingNotification) dto.ThingNotificationExtResponse {
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

func ConvertToThingNotificationsExtResponseDTO(req []models.ExtThingNotification) dto.ThingNotificationsExtResponse {
	res := make([]dto.ThingNotificationExtResponse, 0, len(req))

	for _, notification := range req {
		res = append(res, ConvertToThingNotificationExtResponseDTO(notification))
	}

	return dto.ThingNotificationsExtResponse{Notifications: res}
}
