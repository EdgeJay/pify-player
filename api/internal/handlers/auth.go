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

	if err != nil || cookie == nil {
		log.Println("User not logged in")

		authUrl, err := spotifyService.GetAuthUrl()
		if err != nil {
			return err
		}

		// Instruct client to goto Spotify login page
		return c.JSON(http.StatusOK, loginResponse{LoggedIn: false, RedirectUrl: authUrl})
	} else {
		// check if session is valid
	}

	return c.JSON(http.StatusOK, loginResponse{LoggedIn: true})
}

func getCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "missing_code_or_state"})
	}

	accessToken, err := spotifyService.GetApiToken(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "get_access_token_failed"})
	}

	// get user info
	spotifyUser, err := spotifyService.GetUser(accessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "get_user_info_failed"})
	}

	// save user info into DB
	_, err = userService.SaveUser(spotifyUser.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "save_user_info_failed"})
	}

	// generate UUIDv7 for session ID
	sessionId, err := uuid.NewV7()

	// save session into DB

	// set cookies
	c.SetCookie(utils.CreateCookie(
		COOKIE_SESSION_ID,
		sessionId.String(),
		time.Now().Add(1*time.Hour)),
	)

	return c.Redirect(http.StatusTemporaryRedirect, utils.GetCallbackDestination())
}
