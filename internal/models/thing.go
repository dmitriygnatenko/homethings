package models

import "time"

type Thing struct {
	ID          uint64
	PlaceID     uint64
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
	ID          uint64
	Title       string
	Description string
}

type AddPlaceThingRequest struct {
	PlaceID uint64
	ThingID uint64
}

type UpdatePlaceThingRequest struct {
	ThingID uint64
	PlaceID uint64
}
