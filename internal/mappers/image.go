package mappers

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToImageResponseDTO(req models.Image) dto.ImageResponse {
	var placeIDPtr, thingIDPtr *int

	if req.PlaceID.Valid {
		placeID := int(req.PlaceID.Int64)
		placeIDPtr = &placeID
	}

	if req.ThingID.Valid {
		thingID := int(req.ThingID.Int64)
		thingIDPtr = &thingID
	}

	return dto.ImageResponse{
		ID:        req.ID,
		Image:     req.Image,
		PlaceID:   placeIDPtr,
		ThingID:   thingIDPtr,
		CreatedAt: req.CreatedAt.Format(defaultDateTimeLayout),
	}
}

func ConvertToImagesResponseDTO(images []models.Image) dto.ImagesResponse {
	res := make([]dto.ImageResponse, 0, len(images))

	for _, image := range images {
		res = append(res, ConvertToImageResponseDTO(image))
	}

	return dto.ImagesResponse{Images: res}
}
