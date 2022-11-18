package mappers

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToAddThingRequestModel(req dto.AddThingRequest) models.AddThingRequest {
	res := models.AddThingRequest{
		Title:       req.Title,
		Description: sql.NullString{},
	}

	if req.Description != nil {
		res.Description = sql.NullString{Valid: true, String: *req.Description}
	}

	return res
}

func ConvertToThingResponseDTO(req models.Thing) dto.ThingResponse {
	res := dto.ThingResponse{
		ID:        req.ID,
		Title:     req.Title,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}

	if req.Description.Valid {
		res.Description = req.Description.String
	}

	return res
}
