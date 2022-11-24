package dto

type EmptyResponse struct{}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidateErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
