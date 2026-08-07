package httpx

// ErrorResponse and ValidationErrorResponse document the JSON shapes
// Error() writes at runtime. They exist only so swaggo's @Failure
// annotations have a concrete type to point at; handlers never construct
// them directly (Error() builds the map inline instead).
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" example:"not_found"`
		Message string `json:"message" example:"Task not found"`
	} `json:"error"`
}

type ValidationErrorResponse struct {
	Message string              `json:"message" example:"The given data was invalid."`
	Errors  map[string][]string `json:"errors"`
}
