package models

import (
	"database/sql"
	"time"
)

type Image struct {
	ID        uint64        `db:"id"`
	ThingID   sql.NullInt64 `db:"thing_id"`
	PlaceID   sql.NullInt64 `db:"place_id"`
	Image     string        `db:"image"`
	CreatedAt time.Time     `db:"created_at"`
}
