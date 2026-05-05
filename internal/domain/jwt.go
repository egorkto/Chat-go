package domain

type JWT struct {
	Access  string
	Refresh string
}

func NewJWT(
	access, refresh string,
) JWT {
	return JWT{
		Access:  access,
		Refresh: refresh,
	}
}
