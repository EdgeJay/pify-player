package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/edgejay/pify-player/api/internal/constants"
	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type WSCommand struct {
	Command string
	Payload map[string]string
}

type WSResponse struct {
	Command string            `json:"command"`
	Body    map[string]string `json:"body"`
}

type RemoteResponse struct {
	Success bool `json:"success"`
}

func SetRemoteRoutes(group *echo.Group) {
	group.POST("/play", remoteStartPlayback)
}

func remoteStartPlayback(c echo.Context) error {
	log.Println("remote starting playback...")

	// check if user is already logged in
	cookie, err := c.Cookie(constants.COOKIE_SESSION_ID)
	loggedIn := false

	var session *models.UserSession

	if err == nil && cookie != nil {
		session, err = userService.GetSession(cookie.Value)
		if session != nil && err == nil {
			if time.Now().Before(session.AccessTokenExpiresAt) {
				loggedIn = true
			} else {
				log.Println("access token expired")
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
			if err != nil {
				return c.JSON(http.StatusBadRequest, RemoteResponse{false})
			}

			return c.JSON(http.StatusOK, RemoteResponse{true})
		}
		log.Println("websocket with player not established")
		return c.JSON(http.StatusNoContent, RemoteResponse{false})
	}

	return c.JSON(http.StatusUnauthorized, RemoteResponse{false})
}
