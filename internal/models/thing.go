package models

import "database/sql"

type Thing struct {
	ID          int
	Title       string
	Description sql.NullString
	CreatedAt   string
	UpdatedAt   string
}

type AddThingRequest struct {
	Title       string
	Description sql.NullString
}

type UpdateThingRequest struct {
	PlaceID     sql.NullInt64
	Title       sql.NullString
	Description sql.NullString
}

type AddPlaceThingRequest struct {
	PlaceID int
	ThingID int
}
