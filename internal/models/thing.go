package models

type Thing struct {
	ID          int
	Title       string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

type AddThingRequest struct {
	Title       string
	Description string
}

type UpdateThingRequest struct {
	ID          int
	Title       string
	Description string
}

type AddPlaceThingRequest struct {
	PlaceID int
	ThingID int
}

type UpdatePlaceThingRequest struct {
	ThingID int
	PlaceID int
}
