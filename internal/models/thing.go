package models

import "time"

type Thing struct {
	ID          int
	PlaceID     int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
