package services

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/edgejay/pify-player/api/internal/utils"
)

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyUser struct {
	Id          string `json:"id"`
	DisplayName string `json:"display_name"`
	Images      []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
}

type SpotifyService struct {
	clientId     string
	clientSecret string
	redirectUri  string
	httpClient   *http.Client
}

func GetSpotifyCredentials() SpotifyCredentials {
	return SpotifyCredentials{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
	}
}

func NewSpotifyService(credentials SpotifyCredentials, httpClient *http.Client) *SpotifyService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &SpotifyService{
		clientId:     credentials.ClientID,
		clientSecret: credentials.ClientSecret,
		redirectUri:  credentials.RedirectURI,
		httpClient:   httpClient,
	}
}

func (s *SpotifyService) GetAuthUrl() (string, error) {
	state := utils.GenerateRandomString(16)

	authUrl, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		return "", err
	}

	q := authUrl.Query()
	q.Set("client_id", s.clientId)
	q.Set("response_type", "code")
	q.Set("redirect_uri", s.redirectUri)
	q.Set("state", state)
	q.Set("scope", strings.Join(s.GetScope(), " "))
	authUrl.RawQuery = q.Encode()

	return authUrl.String(), nil
}

func (s *SpotifyService) GetApiToken(code string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectUri)

	tokenReq, err := http.NewRequest(
		"POST",
		"https://accounts.spotify.com/api/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.SetBasicAuth(s.clientId, s.clientSecret)
	tokenRes, err := s.httpClient.Do(tokenReq)
	if err != nil {
		return nil, err
	}
	defer tokenRes.Body.Close()

	tokenResJson := SpotifyTokenResponse{}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenResJson); err != nil {
		return nil, err
	}

	return &tokenResJson, nil
}

func (s *SpotifyService) RefreshApiToken(refreshToken string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", s.clientId)

	tokenReq, err := http.NewRequest(
		"POST",
		"https://accounts.spotify.com/api/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// tokenReq.SetBasicAuth(s.clientId, s.clientSecret)
	tokenRes, err := s.httpClient.Do(tokenReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tokenRes.Body.Close()

	tokenResJson := SpotifyTokenResponse{}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenResJson); err != nil {
		return nil, err
	}

	return &tokenResJson, nil
}

func (s *SpotifyService) GetUser(accessToken string) (*SpotifyUser, error) {
	userReq, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	userReq.Header.Set("Authorization", "Bearer "+accessToken)
	userRes, err := s.httpClient.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer userRes.Body.Close()

	spotifyUser := SpotifyUser{}
	if err := json.NewDecoder(userRes.Body).Decode(&spotifyUser); err != nil {
		return nil, err
	}

	return &spotifyUser, nil
}

func (s *SpotifyService) GetScope() []string {
	return []string{
		"user-read-email",
		"user-read-private",
		"streaming",
		"user-read-playback-state",
		"user-modify-playback-state",
		"user-read-currently-playing",
		"playlist-read-private",
		"playlist-read-collaborative",
		"user-library-read",
		"user-read-playback-position",
		"user-read-recently-played",
		"user-top-read",
	}
}
