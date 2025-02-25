package services

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/edgejay/pify-player/api/internal/utils"
)

type SpotifyService struct {
	clientId    string
	redirectUri string
	httpClient  *http.Client
}

func NewSpotifyService(clientId, redirectUri string, httpClient *http.Client) *SpotifyService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return &SpotifyService{clientId, redirectUri, httpClient}
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
