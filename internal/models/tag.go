package models

import "time"

type Tag struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	Style     string    `db:"style"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AddTagRequest struct {
	Title string
	Style string
}

type UpdateTagRequest struct {
	ID    uint64
	Title string
	Style string
}
