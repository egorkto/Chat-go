package auth_transport

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=20" example:"Иван Иванов"`
	Password string `json:"password" validate:"required,min=5,max=100" example:"Qwer1234?,4"`
}

type SignUpResponse struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token" example:"access_token_example"`
	RefreshToken string  `json:"refresh_token" example:"refresh_token_example"`
}

type UserDTO struct {
	ID       int    `json:"id" example:"1"`
	Version  int    `json:"version" example:"1"`
	FullName string `json:"full_name" example:"Иван Иванов"`
}
