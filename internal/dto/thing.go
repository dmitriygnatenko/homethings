package dto

type AddThingRequest struct {
	PlaceID     uint64 `json:"place_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateThingRequest struct {
	PlaceID     uint64 `json:"place_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type ThingResponse struct {
	ID          uint64 `json:"id"`
	PlaceID     uint64 `json:"place_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ThingsResponse struct {
	Things []ThingResponse `json:"things"`
}

type ThingExtResponse struct {
	ThingResponse
	Tags []TagResponse `json:"tags"`
}

type ThingsExtResponse struct {
	Things []ThingExtResponse `json:"things"`
}
