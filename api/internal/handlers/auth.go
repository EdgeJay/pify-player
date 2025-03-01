package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/edgejay/pify-player/api/internal/utils"
)

type loginResponse struct {
	LoggedIn    bool   `json:"logged_in"`
	RedirectUrl string `json:"redirect_url"`
	ErrorCode   string `json:"error_code"`
}

type callbackPayload struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

const (
	COOKIE_SESSION_ID = "pify_user_sess_id"
)

var spotifyService *services.SpotifyService = services.NewSpotifyService(services.GetSpotifyCredentials(), nil)

var userService *services.UserService = services.NewUserService(database.GetSQLiteDB())

func SetAuthRoutes(group *echo.Group) {
	group.GET("/login", login)
	group.GET("/callback", getCallback)
}

func login(c echo.Context) error {
	// check if user is already logged in
	cookie, err := c.Cookie(COOKIE_SESSION_ID)

	loggedIn := false

	if err == nil && cookie != nil {
		// check if session id stored in cookie is valid
		if session, err := userService.GetSession(cookie.Value); session != nil && err == nil {
			// check accessToken expiry
			if time.Now().After(session.AccessTokenExpiresAt) {
				// refresh access token
				log.Printf("access token expired at %v\n", session.AccessTokenExpiresAt)
				log.Println("refreshing access token...")
				if res, err := spotifyService.RefreshApiToken(session.RefreshToken); err == nil {
					log.Println("got new access token")
					// save access token into DB
					log.Println(session.Uuid)
					if _, err := userService.UpdateSessionAccessToken(
						session.Uuid,
						res.AccessToken,
						time.Now().Add(time.Duration(res.ExpiresIn)*time.Second),
					); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			// valid session
			loggedIn = true
		}
	}

	if !loggedIn {
		log.Println("User not logged in")

		authUrl, err := spotifyService.GetAuthUrl()
		if err != nil {
			return err
		}

		// Instruct client to goto Spotify login page
		return c.JSON(http.StatusOK, loginResponse{LoggedIn: false, RedirectUrl: authUrl})
	} else {
		log.Printf("existing session %s found \n", cookie.Value)
	}

	return c.JSON(http.StatusOK, loginResponse{LoggedIn: true})
}

func getCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "missing_code_or_state"})
	}

	tokenRes, err := spotifyService.GetApiToken(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "get_access_token_failed"})
	}

	// get user info
	spotifyUser, err := spotifyService.GetUser(tokenRes.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "get_user_info_failed"})
	}

	// save user info into DB
	dbUser, err := userService.SaveUser(spotifyUser.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "save_user_info_failed"})
	}

	// generate UUIDv7 for session ID
	sessionId, err := uuid.NewV7()

	// save session into DB
	session, err := userService.SaveSession(
		dbUser.Id,
		sessionId.String(),
		c.Request().UserAgent(),
		tokenRes.AccessToken,
		tokenRes.RefreshToken,
		time.Now().Add(time.Duration(tokenRes.ExpiresIn)*time.Second),
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "save_session_failed"})
	}

	log.Println("user session created:", session.Uuid)

	// set cookies
	c.SetCookie(utils.CreateCookie(
		COOKIE_SESSION_ID,
		sessionId.String(),
		time.Now().Add(720*time.Hour)), // 30 days
	)

	return c.Redirect(http.StatusTemporaryRedirect, utils.GetCallbackDestination())
}
