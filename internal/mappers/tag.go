package mappers

import (
	"sort"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToAddTagRequestModel(req dto.AddTagRequest) models.AddTagRequest {
	return models.AddTagRequest{
		Title: req.Title,
		Style: req.Style,
	}
}

func ConvertToUpdateTagRequestModel(id int, req dto.UpdateTagRequest) models.UpdateTagRequest {
	return models.UpdateTagRequest{
		ID:    id,
		Title: req.Title,
		Style: req.Style,
	}
}

func ConvertToAddThingTagRequestModel(tagID int, thingID int) models.AddThingTagRequest {
	return models.AddThingTagRequest{
		ThingID: thingID,
		TagID:   tagID,
	}
}

func ConvertToDeleteThingTagRequestModel(tagID int, thingID int) models.DeleteThingTagRequest {
	return models.DeleteThingTagRequest{
		ThingID: thingID,
		TagID:   tagID,
	}
}

func ConvertToTagResponseDTO(req models.Tag) dto.TagResponse {
	return dto.TagResponse{
		ID:        req.ID,
		Title:     req.Title,
		Style:     req.Style,
		CreatedAt: req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt: req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ConvertToTagsResponseDTO(req []models.Tag) dto.TagsResponse {
	res := make([]dto.TagResponse, 0, len(req))

	for _, p := range req {
		res = append(res, ConvertToTagResponseDTO(p))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Title < res[j].Title
	})

	return dto.TagsResponse{
		Tags: res,
	}
}
