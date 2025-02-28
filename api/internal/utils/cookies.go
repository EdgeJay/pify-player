package utils

import (
	"net/http"
	"time"
)

func CreateCookie(name, value string, duration time.Time) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = value
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = duration
	return cookie
}
