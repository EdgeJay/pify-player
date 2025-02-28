package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/edgejay/pify-player/api/internal/services"
	"github.com/edgejay/pify-player/api/internal/utils"
)

const (
	COOKIE_SESSION_ID = "pify_user_sess_id"
)

type loginResponse struct {
	LoggedIn    bool   `json:"logged_in"`
	RedirectUrl string `json:"redirect_url"`
	ErrorCode   string `json:"error_code"`
}

type callbackPayload struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyUserResponse struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
	DisplayName string `json:"display_name"`
	Images      []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
}

func SetAuthRoutes(group *echo.Group) {
	group.GET("/login", login)
	group.GET("/callback", getCallback)
}

func login(c echo.Context) error {
	// check if user is already logged in
	cookie, err := c.Cookie(COOKIE_SESSION_ID)

	if err != nil || cookie == nil {
		log.Println("User not logged in")

		spotifyCredentials := utils.GetSpotifyCredentials()

		clientId := spotifyCredentials.ClientID
		redirectUri := spotifyCredentials.RedirectURI

		spotifyService := services.NewSpotifyService(clientId, redirectUri, nil)
		authUrl, err := spotifyService.GetAuthUrl()
		if err != nil {
			return err
		}

		// Instruct client to goto Spotify login page
		return c.JSON(http.StatusOK, loginResponse{LoggedIn: false, RedirectUrl: authUrl})
	} else {
		// check if session is valid
	}

	return c.JSON(http.StatusOK, loginResponse{LoggedIn: true})
}

func getCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, loginResponse{LoggedIn: false, ErrorCode: "missing_code_or_state"})
	}

	// get access token
	spotifyCredentials := utils.GetSpotifyCredentials()
	clientId := spotifyCredentials.ClientID
	clientSecret := spotifyCredentials.ClientSecret
	redirectUri := spotifyCredentials.RedirectURI

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
		log.Println(err)
		return err
	}
	defer tokenRes.Body.Close()

	tokenResJson := SpotifyTokenResponse{}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenResJson); err != nil {
		return err
	}

	// set cookies
	cookie := &http.Cookie{}
	cookie.Name = COOKIE_SESSION_ID
	cookie.Value = tokenResJson.AccessToken
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = time.Now().Add(1 * time.Hour)
	log.Println("cookie is valid:", cookie.Valid())
	c.SetCookie(cookie)

	/*
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
		userResJson.AccessToken = tokenResJson.AccessToken

		return c.JSON(http.StatusOK, userResJson)
	*/

	return c.Redirect(http.StatusTemporaryRedirect, utils.GetCallbackDestination())
}
