package models

import (
	"database/sql"
	"time"
)

type Place struct {
	ID        uint64
	Title     string
	ParentID  sql.NullInt64
	CreatedAt time.Time
	UpdatedAt time.Time
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
