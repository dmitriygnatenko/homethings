package dto

type EmptyResponse struct{}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
