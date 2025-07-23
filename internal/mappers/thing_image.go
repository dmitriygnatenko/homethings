package mappers

import "github.com/dmitriygnatenko/homethings/internal/models"

func ToAddThingImageRequest(thingID uint64, image string) models.AddThingImageRequest {
	return models.AddThingImageRequest{
		ThingID: thingID,
		Image:   image,
	}
}
