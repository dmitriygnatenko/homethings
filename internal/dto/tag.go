package dto

type TagsResponse struct {
	Tags []TagResponse `json:"tags"`
}

type TagResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
