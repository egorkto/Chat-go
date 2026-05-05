package transport_http

import (
	"net/http"
	"time"
)

func NewCookie(
	name, value string,
	expires time.Time,
	path string,
	httpOnly bool,
) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Path:     path,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteLaxMode,
	}
}
