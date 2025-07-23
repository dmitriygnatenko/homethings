package mappers

import "github.com/dmitriygnatenko/homethings/internal/models"

func ToAddPlaceImageRequest(placeID uint64, image string) models.AddPlaceImageRequest {
	return models.AddPlaceImageRequest{
		PlaceID: placeID,
		Image:   image,
	}
}
