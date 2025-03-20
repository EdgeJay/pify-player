package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/edgejay/pify-player/api/internal/errors"
	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/services"
)

func SetDeviceRoutes(group *echo.Group) {
	group.GET("/all", allDevices, middlewareFactory.Auth(), middlewareFactory.GetSpotifyService())
}

func allDevices(c echo.Context) error {
	session := c.Get("session").(*models.UserSession)

	// retrieve all devices
	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	devicesRes, err := spotifyService.GetUserDevices(session.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			Data:      nil,
			ErrorCode: errors.GET_DEVICES_FAILED,
		})
	}
	return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
		Data:      devicesRes,
		ErrorCode: "",
	})
}
