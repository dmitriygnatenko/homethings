package models

type ThingTag struct {
	Tag
	ThingID int
}

type AddThingTagRequest struct {
	ThingID int
	TagID   int
}

type DeleteThingTagRequest struct {
	ThingID int
	TagID   int
}
