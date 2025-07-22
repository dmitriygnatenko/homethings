package mappers

import (
	"sort"

	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func ToAddTagRequest(req dto.AddTagRequest) models.AddTagRequest {
	return models.AddTagRequest{
		Title: req.Title,
		Style: req.Style,
	}
}

func ToUpdateTagRequest(id uint64, req dto.UpdateTagRequest) models.UpdateTagRequest {
	return models.UpdateTagRequest{
		ID:    id,
		Title: req.Title,
		Style: req.Style,
	}
}

func ToAddThingTagRequest(tagID uint64, thingID uint64) models.AddThingTagRequest {
	return models.AddThingTagRequest{
		ThingID: thingID,
		TagID:   tagID,
	}
}

func ToDeleteThingTagRequest(tagID uint64, thingID uint64) models.DeleteThingTagRequest {
	return models.DeleteThingTagRequest{
		ThingID: thingID,
		TagID:   tagID,
	}
}

func ToTagResponse(req models.Tag) dto.TagResponse {
	return dto.TagResponse{
		ID:        req.ID,
		Title:     req.Title,
		Style:     req.Style,
		CreatedAt: req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt: req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ThingTagToTagResponse(req models.ThingTag) dto.TagResponse {
	return dto.TagResponse{
		ID:        req.ID,
		Title:     req.Title,
		Style:     req.Style,
		CreatedAt: req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt: req.UpdatedAt.Format(defaultDateTimeLayout),
	}
}

func ToTagsResponse(req []models.Tag) dto.TagsResponse {
	res := make([]dto.TagResponse, 0, len(req))

	for _, p := range req {
		res = append(res, ToTagResponse(p))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Title < res[j].Title
	})

	return dto.TagsResponse{
		Tags: res,
	}
}
