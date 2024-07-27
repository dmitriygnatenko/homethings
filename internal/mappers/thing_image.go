package mappers

import "git.dmitriygnatenko.ru/dima/homethings/internal/models"

func ToAddThingImageRequest(thingID int, image string) models.AddThingImageRequest {
	return models.AddThingImageRequest{
		ThingID: thingID,
		Image:   image,
	}
}
