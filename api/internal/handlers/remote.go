package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type RemoteResponse struct {
	Success bool `json:"success"`
}

func SetRemoteRoutes(group *echo.Group) {
	group.POST("/play", remoteStartPlayback)
}

func remoteStartPlayback(c echo.Context) error {
	// check if user is already logged in
	cookie, err := c.Cookie(COOKIE_SESSION_ID)
	loggedIn := false

	var session *models.UserSession

	if err == nil && cookie != nil {
		session, err = userService.GetSession(cookie.Value)
		if session != nil && err == nil {
			if time.Now().Before(session.AccessTokenExpiresAt) {
				loggedIn = true
			}
		}
	}

	if loggedIn {
		if playerWebsocket != nil {
			payload, err := json.Marshal(WSResponse{
				Command: "connect",
				Body: map[string]string{
					"access_token": session.AccessToken,
				},
			})
			if err != nil {
				return c.JSON(http.StatusBadRequest, RemoteResponse{false})
			}

			err = websocket.Message.Send(playerWebsocket, string(payload))

			return c.JSON(http.StatusOK, RemoteResponse{true})
		}
		log.Println("websocket with player not established")
		return c.JSON(http.StatusNoContent, RemoteResponse{false})
	}

	return c.JSON(http.StatusUnauthorized, RemoteResponse{false})
}
