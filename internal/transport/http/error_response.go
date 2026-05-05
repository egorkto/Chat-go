package transport_http

type ErrorResponse struct {
	Message string `json:"message" example:"failed to proccess request"`
	Err     string `json:"error" example:"error details"`
}
