package models

import "time"

type Thing struct {
	ID          uint64    `db:"id"`
	PlaceID     uint64    `db:"place_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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
