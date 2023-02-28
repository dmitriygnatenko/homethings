package dto

type AddTagRequest struct {
	Title string `json:"title" validate:"required"`
	Style string `json:"style" validate:"required"`
}

type UpdateTagRequest struct {
	Title string `json:"title" validate:"required"`
	Style string `json:"style" validate:"required"`
}

type TagResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Style     string `json:"style"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
