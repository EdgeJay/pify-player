package middlewares

import (
	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/labstack/echo/v4"
)

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

func (mw *MiddlewareFactory) GetCookie() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, _ := c.Cookie(mw.cookieSessionId)
			c.Set("cookie", cookie)
			return next(c)
		}
	}
}

func (mw *MiddlewareFactory) GetUserService() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("userService", mw.userService)
			return next(c)
		}
	}
}

func (mw *MiddlewareFactory) GetSpotifyService() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("spotifyService", mw.spotifyService)
			return next(c)
		}
	}
}
