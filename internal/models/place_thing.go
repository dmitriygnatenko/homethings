package models

import "time"

type PlaceThing struct {
	PlaceID   uint64
	ThingID   uint64
	CreatedAt time.Time
}
