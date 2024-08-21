package models

import (
	"database/sql"
	"time"
)

type Image struct {
	ID        uint64
	ThingID   sql.NullInt64
	PlaceID   sql.NullInt64
	Image     string
	CreatedAt time.Time
}
