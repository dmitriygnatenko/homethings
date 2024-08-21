package models

type ThingTag struct {
	Tag
	ThingID uint64
}

type AddThingTagRequest struct {
	ThingID uint64
	TagID   uint64
}

type DeleteThingTagRequest struct {
	ThingID uint64
	TagID   uint64
}
