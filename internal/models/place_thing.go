package models

import "time"

type PlaceThing struct {
	PlaceID   int
	ThingID   int
	CreatedAt time.Time
}
