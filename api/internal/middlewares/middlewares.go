package middlewares

import "github.com/edgejay/pify-player/api/internal/services"

type MiddlewareFactory struct {
	cookieSessionId string
	userService     *services.UserService
	spotifyService  *services.SpotifyService
}

func NewMiddlewareFactory(
	cookieSessionId string,
	userService *services.UserService,
	spotifyService *services.SpotifyService,
) *MiddlewareFactory {
	return &MiddlewareFactory{
		cookieSessionId,
		userService,
		spotifyService,
	}
}
