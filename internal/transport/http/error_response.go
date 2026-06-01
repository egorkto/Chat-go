package transport_http

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

type ValidationErrorResponse struct {
	ErrorResponse
	Details map[string]string `json:"details" example:"field:validation error message"`
}
