package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/services"
)

type loginResponse struct {
	LoggedIn  bool   `json:"logged_in"`
	ErrorCode string `json:"error_code"`
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyUserResponse struct {
	Id          string `json:"id"`
	DisplayName string `json:"display_name"`
	Images      []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
}

func SetAuthRoutes(group *echo.Group) {
	group.GET("/login", login)
	group.GET("/callback", callback)
}

func login(c echo.Context) error {
	// check if user is already logged in
	cookie, err := c.Cookie("pify_user_sess_id")

	if err != nil || cookie == nil {
		clientId := os.Getenv("SPOTIFY_CLIENT_ID")
		redirectUri := os.Getenv("SPOTIFY_REDIRECT_URI")

		spotifyService := services.NewSpotifyService(clientId, redirectUri, nil)
		authUrl, err := spotifyService.GetAuthUrl()
		if err != nil {
			return err
		}

		// Redirect to Spotify login page
		return c.Redirect(http.StatusTemporaryRedirect, authUrl)
	}

	// TODO: check if session is valid

	return c.JSON(http.StatusOK, loginResponse{LoggedIn: true})
}

func callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorCode := c.QueryParam("error")

	if errorCode != "" {
		return c.JSON(http.StatusUnauthorized, loginResponse{LoggedIn: false, ErrorCode: errorCode})
	}

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "missing_code_or_state"})
	}

	// get access token
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectUri := os.Getenv("SPOTIFY_REDIRECT_URI")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)

	tokenReq, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.SetBasicAuth(clientId, clientSecret)

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	tokenRes, err := client.Do(tokenReq)
	if err != nil {
		return err
	}
	defer tokenRes.Body.Close()

	tokenResJson := SpotifyTokenResponse{}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenResJson); err != nil {
		return err
	}

	// get user info
	userReq, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResJson.AccessToken)
	userRes, err := client.Do(userReq)
	if err != nil {
		return err
	}
	defer userRes.Body.Close()

	userResJson := SpotifyUserResponse{}
	if err := json.NewDecoder(userRes.Body).Decode(&userResJson); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userResJson)
}
