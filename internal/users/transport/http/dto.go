package users_transport_http

import (
	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
)

type UserDTOResponse struct {
	ID       int    `json:"id" example:"1"`
	Version  int    `json:"version" example:"1"`
	FullName string `json:"full_name" example:"John Doe"`
	Login    string `json:"login" example:"johndoe"`
}

func domainToDTO(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID(),
		Version:  user.Version(),
		FullName: user.FullName(),
		Login:    user.Login(),
	}
}

type ErrorResponse = transport_http.ErrorResponse
type ValidationErrorResponse = transport_http.ValidationErrorResponse
