package mappers

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToPlaceResponseDTO(req models.Place) dto.PlaceResponse {
	res := dto.PlaceResponse{
		ID:        req.ID,
		Title:     req.Title,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}

	if req.ParentID.Valid {
		parentID := int(req.ParentID.Int64)
		res.ParentID = &parentID
	}

	return res
}

func ConvertToAddPlaceRequestModel(req dto.AddPlaceRequest) models.AddPlaceRequest {
	res := models.AddPlaceRequest{
		Title: req.Title,
	}

	if req.ParentID != nil {
		res.ParentID = sql.NullInt64{Int64: int64(*req.ParentID), Valid: true}
	}

	return res
}

func ConvertToUpdatePlaceRequestModel(id int, req dto.UpdatePlaceRequest) models.UpdatePlaceRequest {
	res := models.UpdatePlaceRequest{
		ID:    id,
		Title: req.Title,
	}

	if req.ParentID != nil {
		res.ParentID = sql.NullInt64{Int64: int64(*req.ParentID), Valid: true}
	}

	return res
}
