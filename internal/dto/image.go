package dto

type ImageResponse struct {
	ID        int    `json:"id"`
	Image     string `json:"image"`
	PlaceID   *int   `json:"place_id"`
	ThingID   *int   `json:"thing_id"`
	CreatedAt string `json:"created_at"`
}

type ImagesResponse struct {
	Images []ImageResponse `json:"images"`
}
