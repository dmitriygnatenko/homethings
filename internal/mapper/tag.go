package mapper

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToTagsResponseDTO(tags []models.Tag) dto.TagsResponse {
	res := make([]dto.TagResponse, 0, len(tags))

	for i := range tags {
		res = append(res, ConvertToTagResponseDTO(tags[i]))
	}

	return dto.TagsResponse{Tags: res}
}

func ConvertToTagResponseDTO(tag models.Tag) dto.TagResponse {
	return dto.TagResponse{
		ID:        tag.ID,
		Title:     tag.Title,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}
