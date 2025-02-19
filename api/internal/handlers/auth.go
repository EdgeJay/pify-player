package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetAuthRoutes(group *echo.Group) {
	group.GET("/hello", hello)
}

// Get all users
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello!"})
}
