package dto

type AddThingRequest struct {
	PlaceID     int    `json:"place_id" validate:"required"`
	Title       string `json:"title"  validate:"required"`
	Description string `json:"description"`
}

type UpdateThingRequest struct {
	PlaceID     *int    `json:"place_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type ThingResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ThingsResponse struct {
	Things []ThingResponse `json:"things"`
}
