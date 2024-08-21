package dto

type AddPlaceRequest struct {
	ParentID *uint64 `json:"parent_id"`
	Title    string  `json:"title" validate:"required"`
}

type UpdatePlaceRequest struct {
	ParentID *uint64 `json:"parent_id"`
	Title    string  `json:"title" validate:"required"`
}

type PlaceResponse struct {
	ID        uint64  `json:"id"`
	ParentID  *uint64 `json:"parent_id"`
	Title     string  `json:"title"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type PlacesResponse struct {
	Places []PlaceResponse `json:"places"`
}

type PlaceTreeResponse struct {
	Places []PlaceTree `json:"places"`
}

type PlaceTree struct {
	Place        PlaceResponse `json:"place"`
	NestedPlaces []PlaceTree   `json:"nested"`
}
