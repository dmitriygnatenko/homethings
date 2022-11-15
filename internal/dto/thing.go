package dto

type AddThingRequest struct {
	PlaceID     int    `json:"place_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateThingRequest struct {
	PlaceID     *int    `json:"place_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type ThingResponse struct {
	ID          int    `json:"id"`
	PlaceID     int    `json:"place_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
