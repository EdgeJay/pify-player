package handlers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func SetAuthRoutes(group *echo.Group) {
	group.GET("/hello", hello)
}

// Get all users
func hello(c echo.Context) error {
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	return c.JSON(http.StatusOK, map[string]string{
		"message":      "Hello!",
		"clientId":     clientId,
		"clientSecret": clientSecret,
	})
}
