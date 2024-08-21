package mappers

import "git.dmitriygnatenko.ru/dima/homethings/internal/models"

func ToAddPlaceImageRequest(placeID uint64, image string) models.AddPlaceImageRequest {
	return models.AddPlaceImageRequest{
		PlaceID: placeID,
		Image:   image,
	}
}
