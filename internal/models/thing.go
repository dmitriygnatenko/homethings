package models

import "database/sql"

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
	Title       sql.NullString
	Description sql.NullString
}

type AddPlaceThingRequest struct {
	PlaceID int
	ThingID int
}

type UpdatePlaceThingRequest struct {
	ThingID int
	PlaceID int
}
