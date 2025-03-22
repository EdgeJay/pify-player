package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/edgejay/pify-player/api/internal/errors"
	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/services"
)

type ControlPlaybackRequest struct {
	AccessToken string `json:"access_token"`
	DeviceId    string `json:"device_id"`
}

func SetDeviceRoutes(group *echo.Group) {
	group.GET("/all", allDevices, middlewareFactory.Auth(), middlewareFactory.GetSpotifyService())
	// following endpoint is only meant to be called from player page only
	group.POST("/control-playback", controlPlayback, middlewareFactory.GetSpotifyService())
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

func controlPlayback(c echo.Context) error {
	var req ControlPlaybackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			Data:      nil,
			ErrorCode: errors.INVALID_REQUEST_BODY,
		})
	}

	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	success, err := spotifyService.TransferPlayback(req.AccessToken, req.DeviceId)
	if success {
		return c.JSON(http.StatusNoContent, nil)
	} else {
		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			Data:      nil,
			ErrorCode: err.Error(),
		})
	}
}
