package auth_transport

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type SignUpResponse struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Version  int    `json:"version"`
	FullName string `json:"full_name"`
}
