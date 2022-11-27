package dto

type AddPlaceRequest struct {
	ParentID *int   `json:"parent_id"`
	Title    string `json:"title" validate:"required"`
}

type UpdatePlaceRequest struct {
	ParentID *int   `json:"parent_id"`
	Title    string `json:"title" validate:"required"`
}

type PlaceResponse struct {
	ID        int    `json:"id"`
	ParentID  *int   `json:"parent_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
