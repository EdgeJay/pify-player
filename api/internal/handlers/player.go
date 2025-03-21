package handlers

import (
	"encoding/json"
	"log"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/services"
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

var playerService *services.PlayerService = services.NewPlayerService(database.GetSQLiteDB())

var playerWebsocket *websocket.Conn

func parseWebsocketMessage(c echo.Context, ws *websocket.Conn, msg string) error {

	command := WSCommand{}

	if err := json.Unmarshal([]byte(msg), &command); err == nil {
		spotifyService := c.Get("spotifyService").(*services.SpotifyService)

		// update database to indicate waiting state
		switch command.Command {
		case "connect":
			playerWebsocket = ws
			accessToken := ""
			session, err := playerService.Connect()
			if err != nil {
				log.Println("websocket connect error:", err)
				return err
			} else {
				accessToken = session.AccessToken
				// check if access token is still valid
				if res, err := spotifyService.CheckAndRefreshApiToken(session.AccessTokenExpiresAt, session.RefreshToken); err != nil {
					return err
				} else if res != nil {
					accessToken = res.AccessToken
				}

				// return access token
				payload, err := json.Marshal(WSResponse{
					Command: "connect",
					Body: map[string]string{
						"access_token": accessToken,
					},
				})

				if err != nil {
					log.Println("unable to marshal websocket response", err)
					return err
				}

				err = websocket.Message.Send(playerWebsocket, string(payload))
				if err != nil {
					log.Println("unable to send websocket response", err)
					return err
				}
			}
		}
	} else {
		log.Println("websocket parse message error:", err)
		return err
	}
	return nil
}

func SetPlayerRoutes(group *echo.Group) {
	group.GET("/ws", playerWebsocketEndpoint, middlewareFactory.GetSpotifyService())
}

func playerWebsocketEndpoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err == nil {
				parseWebsocketMessage(c, ws, msg)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
