package mappers

import (
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
	return models.UpdateThingRequest{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	}
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
		PlaceID:     req.PlaceID,
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
