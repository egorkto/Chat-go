package domain

import "time"

type JWT struct {
	Access         string
	Refresh        string
	RefreshExpires time.Duration
}

func NewJWT(
	access, refresh string,
	refreshExpires time.Duration,
) JWT {
	return JWT{
		Access:         access,
		Refresh:        refresh,
		RefreshExpires: refreshExpires,
	}
}
