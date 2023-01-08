package mappers

import "git.dmitriygnatenko.ru/dima/homethings/internal/models"

func ConvertToAddPlaceImageRequestModel(placeID int, image string) models.AddPlaceImageRequest {
	return models.AddPlaceImageRequest{
		PlaceID: placeID,
		Image:   image,
	}
}
