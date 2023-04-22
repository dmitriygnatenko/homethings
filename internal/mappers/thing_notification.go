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
		UpdatedAt:        req.CreatedAt.Format(defaultDateTimeLayout),
	}
}
