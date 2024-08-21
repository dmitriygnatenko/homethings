package models

import "time"

type Tag struct {
	ID        uint64
	Title     string
	Style     string
	CreatedAt time.Time
	UpdatedAt time.Time
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
