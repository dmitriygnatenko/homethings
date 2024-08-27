package models

import "time"

type PlaceThing struct {
	PlaceID   uint64    `db:"place_id"`
	ThingID   uint64    `db:"thing_id"`
	CreatedAt time.Time `db:"created_at"`
}
