package mappers

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToAddThingRequestModel(req dto.AddThingRequest) models.AddThingRequest {
	return models.AddThingRequest{
		Title:       req.Title,
		Description: req.Description,
	}
}

func ConvertToUpdateThingRequestModel(id int, req dto.UpdateThingRequest) models.UpdateThingRequest {
	res := models.UpdateThingRequest{ID: id}

	if req.Title != nil {
		res.Title = sql.NullString{String: *req.Title, Valid: true}
	}

	if req.Description != nil {
		res.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	return res
}

func ConvertToAddPlaceThingRequestModel(thingID int, placeID int) models.AddPlaceThingRequest {
	return models.AddPlaceThingRequest{
		PlaceID: placeID,
		ThingID: thingID,
	}
}

func ConvertToUpdatePlaceThingRequestModel(thingID int, placeID int) models.UpdatePlaceThingRequest {
	return models.UpdatePlaceThingRequest{
		ThingID: thingID,
		PlaceID: placeID,
	}
}

func ConvertToThingResponseDTO(req models.Thing) dto.ThingResponse {
	return dto.ThingResponse{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}
}

func ConvertToThingsResponseDTO(things []models.Thing) dto.ThingsResponse {
	res := make([]dto.ThingResponse, 0, len(things))

	for _, thing := range things {
		res = append(res, ConvertToThingResponseDTO(thing))
	}

	return dto.ThingsResponse{Things: res}
}
