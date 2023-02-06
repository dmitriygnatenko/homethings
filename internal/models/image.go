package models

import (
	"database/sql"
	"time"
)

type Image struct {
	ID        int
	ThingID   sql.NullInt64
	PlaceID   sql.NullInt64
	Image     string
	CreatedAt time.Time
}
