package auth_jwt_transport_http

import "github.com/egorkto/Chat-go/internal/domain"

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100" example:"Иван Иванов"`
	Login    string `json:"login" validate:"required,min=3,max=25" example:"ivan"`
	Password string `json:"password" validate:"required,min=5,max=100" example:"KGorgsoroee3235oOSNG?>,frgs3"`
}

type LogInRequest struct {
	Login    string `json:"login" validate:"required,max=25" example:"ivan"`
	Password string `json:"password" validate:"required,max=150" example:"KGorgsoroee3235oOSNG?>,frgs3"`
}

type AuthResponse struct {
	User        UserDTO `json:"user"`
	AccessToken string  `json:"access_token" example:"access_token_example"`
}

type UserDTO struct {
	ID       int    `json:"id" example:"1"`
	Version  int    `json:"version" example:"1"`
	FullName string `json:"full_name" example:"Иван Иванов"`
}

func responseFromDomain(u domain.User, accessToken string) AuthResponse {
	return AuthResponse{
		User: UserDTO{
			u.ID(),
			u.Version(),
			u.FullName(),
		},
		AccessToken: accessToken,
	}
}
