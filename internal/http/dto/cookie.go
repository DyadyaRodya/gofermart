package dto

import (
	"net/http"
	"time"
)

const (
	CookieName = "Auth"
	TTL        = time.Hour * 0
	Path       = "/api/"
)

func NewAuthCookie(token, cookieName, path string, ttl time.Duration) *http.Cookie {
	c := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Path:     path,
	}
	if ttl > 0 { // allow to be infinite if ttl == 0
		c.Expires = time.Now().Add(ttl)
	}
	return c
}
