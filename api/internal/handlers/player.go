package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/edgejay/pify-player/api/internal/database"
	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var playerService *services.PlayerService = services.NewPlayerService(database.GetSQLiteDB())

var playerWebsocket *websocket.Conn

func SetPlayerRoutes(group *echo.Group) {
	group.GET("/ws", playerWebsocketEndpoint, middlewareFactory.GetSpotifyService())
	group.GET("/connect", connect, middlewareFactory.GetSpotifyService(), middlewareFactory.BasicAuth())
}

func connect(c echo.Context) error {
	accessToken := ""
	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	session, err := playerService.Connect()
	if err != nil {
		log.Println("playerService connect error:", err)
		return err
	}

	accessToken = session.AccessToken
	// check if access token is still valid
	if res, err := spotifyService.CheckAndRefreshApiToken(session.AccessTokenExpiresAt, session.RefreshToken); err != nil {
		return err
	} else if res != nil {
		accessToken = res.AccessToken
	}

	// return access token
	return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
		Data: map[string]string{
			"access_token": accessToken,
		},
	})
}

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
