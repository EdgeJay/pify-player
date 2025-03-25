package middlewares

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/utils"
)

/*
RequireAuth middleware checks if the user is authenticated by:

1. Looking for a session cookie in the request

2. Validating the session ID from the cookie against the database

3. If session is valid, adds user context to the request

4. If session is invalid/missing, returns 401 Unauthorized

Usage:

	router.GET("/protected-route", handler, Auth())
*/
func (mw *MiddlewareFactory) Auth() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get auth url if needed to redirect user to login page
			authUrl, err := mw.spotifyService.GetAuthUrl()
			if err != nil {
				return err
			}

			cookie, err := c.Cookie(mw.cookieSessionId)

			// error in fetching cookie or it does not present
			if err != nil || cookie == nil {
				// Instruct client to goto Spotify login page
				return c.JSON(http.StatusOK, pifyHttp.LoginResponse{LoggedIn: false, RedirectUrl: authUrl})
			}

			session, err := mw.userService.GetSession(cookie.Value)
			if err != nil || session == nil {
				// Instruct client to goto Spotify login page
				return c.JSON(http.StatusOK, pifyHttp.LoginResponse{LoggedIn: false, RedirectUrl: authUrl})
			}

			if res, err := mw.spotifyService.CheckAndRefreshApiToken(session.AccessTokenExpiresAt, session.RefreshToken); err != nil {
				// If token can't be refreshhed, instruct client to goto Spotify login page
				return c.JSON(http.StatusOK, pifyHttp.LoginResponse{LoggedIn: false, RedirectUrl: authUrl})
			} else if res != nil {
				log.Println("got new access token:", res.AccessToken)
				log.Println(session.Uuid)
				// save access token into DB
				if _, err := mw.userService.UpdateSessionAccessToken(
					session.Uuid,
					res.AccessToken,
					time.Now().Add(time.Duration(res.ExpiresIn)*time.Second),
				); err != nil {
					return err
				}
			} else {
				// no updates needed, access token is still valid
			}

			// Get session again
			session, err = mw.userService.GetSession(cookie.Value)
			if err != nil || session == nil {
				// Instruct client to goto Spotify login page
				return c.JSON(http.StatusOK, pifyHttp.LoginResponse{LoggedIn: false, RedirectUrl: authUrl})
			}

			// Add session to context for downstream handlers
			c.Set("session", session)

			return next(c)
		}
	}
}

// BasicAuth creates a middleware that performs basic authentication
func (mw *MiddlewareFactory) BasicAuth() func(echo.HandlerFunc) echo.HandlerFunc {
	username := utils.GetBasicAuthUsername()
	password := utils.GetBasicAuthPassword()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract credentials from the request header
			reqUsername, reqPassword, ok := c.Request().BasicAuth()
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
			}

			// Timing attack safe comparison
			if subtle.ConstantTimeCompare([]byte(username), []byte(reqUsername)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(reqPassword)) == 1 {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
	}
}
