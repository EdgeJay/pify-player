package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/errors"
	pifyHttp "github.com/edgejay/pify-player/api/internal/http"
	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/edgejay/pify-player/api/internal/utils"
	"github.com/labstack/echo/v4"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var playerService *services.PlayerService = services.NewPlayerService(database.GetSQLiteDB())

func SetPlayerRoutes(group *echo.Group) {
	group.GET("/connect", connect, middlewareFactory.GetUserService(), middlewareFactory.GetSpotifyService(), middlewareFactory.BasicAuth())
	group.GET("/track/:id", getTrack, middlewareFactory.GetSpotifyService(), middlewareFactory.BasicAuth())
	group.POST("/youtube", getAndSaveYoutubeVideo, middlewareFactory.GetSpotifyService(), middlewareFactory.BasicAuth())
}

func connect(c echo.Context) error {
	accessToken := ""
	var expiresAt time.Time

	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	userService := c.Get("userService").(*services.UserService)

	session, err := playerService.Connect()
	if err != nil {
		if err.Error() == errors.INVALID_SESSION {
			return c.JSON(http.StatusUnauthorized, pifyHttp.ApiResponse{
				ErrorCode: errors.INVALID_SESSION,
			})
		}
		log.Println("playerService connect error:", err)
		return err
	}

	accessToken = session.AccessToken
	expiresAt = session.AccessTokenExpiresAt
	// check if access token is still valid
	if res, err := spotifyService.CheckAndRefreshApiToken(session.AccessTokenExpiresAt, session.RefreshToken); err != nil {
		return err
	} else if res != nil {
		// save access token into DB
		if _, err := userService.UpdateSessionAccessToken(
			session.Uuid,
			res.AccessToken,
			time.Now().Add(time.Duration(res.ExpiresIn)*time.Second),
		); err != nil {
			return err
		}

		accessToken = res.AccessToken
		expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	}

	// return access token
	return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
		Data: map[string]string{
			"access_token": accessToken,
			"expires_at":   expiresAt.Format(time.RFC3339),
		},
	})
}

func getTrack(c echo.Context) error {
	trackId := c.Param("id")
	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	session, err := playerService.Connect()
	if err != nil {
		log.Println("playerService connect error:", err)
		return err
	}

	// check if access token is still valid
	if spotifyService.IsApiTokenExpired(session.AccessTokenExpiresAt); err != nil {
		return c.JSON(http.StatusUnauthorized, pifyHttp.ApiResponse{
			ErrorCode: errors.BAD_OR_EXPIRED_TOKEN,
		})
	}

	track, err := spotifyService.GetTrackBytes(session.AccessToken, trackId)
	if err != nil || track == nil {
		if err.Error() == errors.BAD_OR_EXPIRED_TOKEN {
			return c.JSON(http.StatusUnauthorized, pifyHttp.ApiResponse{
				ErrorCode: errors.BAD_OR_EXPIRED_TOKEN,
			})
		}

		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			ErrorCode: errors.GET_TRACK_FAILED,
		})
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(track, &data); err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			ErrorCode: errors.PARSE_TRACK_RESPONSE_FAILED,
		})
	}

	return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
		Data: data,
	})
}

func getAndSaveYoutubeVideo(c echo.Context) error {
	spotifyService := c.Get("spotifyService").(*services.SpotifyService)
	session, err := playerService.Connect()
	if err != nil {
		log.Println("playerService connect error:", err)
		return err
	}

	// check if access token is still valid
	if spotifyService.IsApiTokenExpired(session.AccessTokenExpiresAt); err != nil {
		return c.JSON(http.StatusUnauthorized, pifyHttp.ApiResponse{
			ErrorCode: errors.BAD_OR_EXPIRED_TOKEN,
		})
	}

	var vidReq pifyHttp.YoutubeVideoRequest
	if err := c.Bind(&vidReq); err != nil {
		return c.JSON(http.StatusBadRequest, pifyHttp.ApiResponse{
			ErrorCode: errors.INVALID_REQUEST_BODY,
		})
	}

	// check if cached result for youtube video exists
	trackMedia := playerService.GetTrackMedia(vidReq.SpotifyTrackId, services.TRACK_MEDIA_TYPE_YOUTUBE)
	if trackMedia != nil {
		// return cached result
		log.Println("Found cached youtube video id:", trackMedia.MediaId)
		return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
			Data: pifyHttp.YoutubeVideoResponse{
				VideoId: trackMedia.MediaId,
			},
		})
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(utils.GetYoutubeApiKey()))
	if err != nil {
		log.Println("Error creating YouTube client: %v", err)
		return err
	}

	call := service.Search.List([]string{"snippet"}).
		Q(vidReq.Query).
		Type("video").
		MaxResults(1)

	serverSettings := utils.GetServerSettings()
	call.Header().Set("Referer", "https://"+serverSettings.SslDomain)

	response, err := call.Do()
	if err != nil {
		log.Println("Error making search API call: %v", err)
		return err
	}

	if len(response.Items) == 0 {
		return c.JSON(http.StatusNotFound, pifyHttp.ApiResponse{
			ErrorCode: errors.NO_YOUTUBE_VIDEO_FOUND,
		})
	}

	videoId := response.Items[0].Id.VideoId

	// save youtube video id to DB
	playerService.SaveTrackMedia(vidReq.SpotifyTrackId, videoId, services.TRACK_MEDIA_TYPE_YOUTUBE)

	return c.JSON(http.StatusOK, pifyHttp.ApiResponse{
		Data: pifyHttp.YoutubeVideoResponse{
			VideoId: videoId,
		},
	})
}
