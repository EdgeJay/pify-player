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

func parseWebsocketMessage(ws *websocket.Conn, msg string) error {
	command := WSCommand{}
	if err := json.Unmarshal([]byte(msg), &command); err == nil {
		// update database to indicate waiting state
		switch command.Command {
		case "connect":
			playerWebsocket = ws
			err := playerService.SetWaitingState()
			if err != nil {
				log.Println("websocket connect error:", err)
				return err
			}
		}
	} else {
		log.Println("websocket parse message error:", err)
		return err
	}
	return nil
}

func SetPlayerRoutes(group *echo.Group) {
	group.GET("/ws", playerWebsocketEndpoint)
}

func playerWebsocketEndpoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err == nil {
				log.Println(parseWebsocketMessage(ws, msg))
			} else {
				// c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
