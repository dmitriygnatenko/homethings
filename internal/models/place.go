package models

import "database/sql"

type Place struct {
	ID        int
	Title     string
	ParentID  sql.NullInt64
	CreatedAt string
	UpdatedAt string
}

type AddPlaceRequest struct {
	Title    string
	ParentID sql.NullInt64
}

type UpdatePlaceRequest struct {
	ID       int
	Title    string
	ParentID sql.NullInt64
}
