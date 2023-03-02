package models

import "time"

type Tag struct {
	ID        int
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
	ID    int
	Title string
	Style string
}
