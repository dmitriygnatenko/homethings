package mappers

import (
	"sort"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ToAddThingRequest(req dto.AddThingRequest) models.AddThingRequest {
	return models.AddThingRequest{
		Title:       req.Title,
		Description: req.Description,
	}
}

func ToUpdateThingRequest(id uint64, req dto.UpdateThingRequest) models.UpdateThingRequest {
	return models.UpdateThingRequest{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	}
}

func ToAddPlaceThingRequest(thingID uint64, placeID uint64) models.AddPlaceThingRequest {
	return models.AddPlaceThingRequest{
		PlaceID: placeID,
		ThingID: thingID,
	}
}

func ToUpdatePlaceThingRequest(thingID uint64, placeID uint64) models.UpdatePlaceThingRequest {
	return models.UpdatePlaceThingRequest{
		ThingID: thingID,
		PlaceID: placeID,
	}
}

func ToThingResponse(req models.Thing) dto.ThingResponse {
	return dto.ThingResponse{
		ID:          req.ID,
		PlaceID:     req.PlaceID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt:   req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ToThingsResponse(things []models.Thing) dto.ThingsResponse {
	res := make([]dto.ThingResponse, 0, len(things))

	for _, thing := range things {
		res = append(res, ToThingResponse(thing))
	}

	return dto.ThingsResponse{Things: res}
}

func ToThingsExtResponse(things []models.Thing, tags []models.ThingTag) dto.ThingsExtResponse {
	res := make([]dto.ThingExtResponse, 0, len(things))

	groupedTags := make(map[uint64][]dto.TagResponse)
	for _, tag := range tags {
		if _, ok := groupedTags[tag.ThingID]; !ok {
			groupedTags[tag.ThingID] = make([]dto.TagResponse, 0)
		}

		groupedTags[tag.ThingID] = append(groupedTags[tag.ThingID], ThingTagToTagResponse(tag))
	}

	for _, thing := range things {
		var thingTags []dto.TagResponse
		if t, ok := groupedTags[thing.ID]; ok {
			sort.Slice(t, func(i, j int) bool {
				return t[i].Title < t[j].Title
			})
			thingTags = t
		}

		res = append(res, dto.ThingExtResponse{
			ThingResponse: ToThingResponse(thing),
			Tags:          thingTags,
		})
	}

	return dto.ThingsExtResponse{Things: res}
}
