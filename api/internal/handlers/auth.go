package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/constants"
	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/edgejay/pify-player/api/internal/errors"
	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/utils"
)

func SetAuthRoutes(group *echo.Group) {
	group.GET("/login", login, middlewareFactory.Auth())
	group.GET("/callback", getCallback)
}

func login(c echo.Context) error {
	session := c.Get("session").(*models.UserSession)
	log.Printf("existing session %s found \n", session.Uuid)
	return c.JSON(http.StatusOK, pifyHttp.LoginResponse{LoggedIn: true})
}

func getCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.MISSING_CODE_OR_STATE})
	}

	tokenRes, err := spotifyService.GetApiToken(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.GET_ACCESS_TOKEN_FAILED})
	}

	// get user info
	spotifyUser, err := spotifyService.GetUser(tokenRes.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.GET_USER_INFO_FAILED})
	}

	// save user info into DB
	dbUser, err := userService.SaveUser(spotifyUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.SAVE_USER_INFO_FAILED})
	}

	// generate UUIDv7 for session ID
	sessionId, err := uuid.NewV7()
	if err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.GENERATE_SESSION_ID_FAILED})
	}

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
		return c.JSON(http.StatusBadRequest, pifyHttp.LoginResponse{LoggedIn: false, ErrorCode: errors.SAVE_SESSION_FAILED})
	}

	log.Println("user session created:", session.Uuid)

	// set cookies
	c.SetCookie(utils.CreateCookie(
		constants.COOKIE_SESSION_ID,
		sessionId.String(),
		time.Now().Add(720*time.Hour)), // 30 days
	)

	return c.Redirect(http.StatusTemporaryRedirect, utils.GetCallbackDestination())
}
