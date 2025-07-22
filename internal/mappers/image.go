package mappers

import (
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func ToImageResponse(req models.Image) dto.ImageResponse {
	var placeIDPtr, thingIDPtr *uint64

	if req.PlaceID.Valid {
		placeID := uint64(req.PlaceID.Int64)
		placeIDPtr = &placeID
	}

	if req.ThingID.Valid {
		thingID := uint64(req.ThingID.Int64)
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

func ToImagesResponse(images []models.Image) dto.ImagesResponse {
	res := make([]dto.ImageResponse, 0, len(images))

	for _, image := range images {
		res = append(res, ToImageResponse(image))
	}

	return dto.ImagesResponse{Images: res}
}
