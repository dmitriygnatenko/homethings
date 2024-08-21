package dto

type ImageResponse struct {
	ID        uint64  `json:"id"`
	Image     string  `json:"image"`
	PlaceID   *uint64 `json:"place_id"`
	ThingID   *uint64 `json:"thing_id"`
	CreatedAt string  `json:"created_at"`
}

type ImagesResponse struct {
	Images []ImageResponse `json:"images"`
}
