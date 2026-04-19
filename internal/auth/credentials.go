package auth

type Credentials struct {
	Password string
}

func NewCredentials(password string) Credentials {
	return Credentials{
		Password: password,
	}
}
