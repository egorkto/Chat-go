package users_transport_http

import "github.com/egorkto/Chat-go/internal/domain"

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Version  int    `json:"version"`
	FullName string `json:"full_name"`
	Login    string `json:"login"`
}

func domainToDTO(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID(),
		Version:  user.Version(),
		FullName: user.FullName(),
		Login:    user.Login(),
	}
}
