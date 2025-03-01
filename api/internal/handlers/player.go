package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type WSCommand struct {
	Command string
	Payload map[string]string
}

var playerService *services.PlayerService = services.NewPlayerService(database.GetSQLiteDB())

func parseWebsocketMessage(msg string) error {
	command := WSCommand{}
	if err := json.Unmarshal([]byte(msg), &command); err == nil {
		// update database to indicate waiting state
		switch command.Command {
		case "connect":
			err := playerService.SetWaitingState()
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
		return err
	}
	return nil
}

func SetPlayerRoutes(group *echo.Group) {
	group.GET("/ws", playerWebSocket)
}

func playerWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			/*
				// Write
				err := websocket.Message.Send(ws, "Hello, Client!")
				if err != nil {
					c.Logger().Error(err)
				}
			*/

			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err == nil {
				fmt.Println(parseWebsocketMessage(msg))
			} else {
				// c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
