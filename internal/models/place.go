package models

import (
	"database/sql"
	"time"
)

type Place struct {
	ID        uint64        `db:"id"`
	Title     string        `db:"title"`
	ParentID  sql.NullInt64 `db:"parent_id"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

type AddPlaceRequest struct {
	Title    string
	ParentID sql.NullInt64
}

type UpdatePlaceRequest struct {
	ID       uint64
	Title    string
	ParentID sql.NullInt64
}
