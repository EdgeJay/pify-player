package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	pifyErrors "github.com/edgejay/pify-player/api/internal/errors"
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

type SpotifyDevice struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
	SupportsVolume   bool   `json:"supports_volume"`
}

type SpotifyDevices struct {
	Devices []SpotifyDevice `json:"devices"`
}

type TransferPlaybackRequest struct {
	DeviceIds []string `json:"device_ids"`
	Play      bool     `json:"play"`
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

func (s *SpotifyService) CheckAndRefreshApiToken(accessTokenExpiresAt time.Time, refreshToken string) (*SpotifyTokenResponse, error) {
	// check accessToken expiry
	if time.Now().After(accessTokenExpiresAt) {
		// refresh access token
		log.Printf("access token expired at %v\n", accessTokenExpiresAt)
		log.Println("refreshing access token...")
		return s.RefreshApiToken(refreshToken)
	}
	return nil, nil
}

func (s *SpotifyService) RefreshApiToken(refreshToken string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

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

func (s *SpotifyService) GetUserDevices(accessToken string) (*SpotifyDevices, error) {
	deviceReq, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return nil, err
	}

	deviceReq.Header.Set("Authorization", "Bearer "+accessToken)
	deviceRes, err := s.httpClient.Do(deviceReq)
	if err != nil {
		return nil, err
	}
	defer deviceRes.Body.Close()

	spotifyDevices := SpotifyDevices{}
	if err := json.NewDecoder(deviceRes.Body).Decode(&spotifyDevices); err != nil {
		return nil, err
	}

	return &spotifyDevices, nil
}

func (s *SpotifyService) TransferPlayback(accessToken, deviceId string) (bool, error) {
	reqPayload := TransferPlaybackRequest{
		DeviceIds: []string{deviceId},
		Play:      false,
	}
	b, err := json.Marshal(reqPayload)
	if err != nil {
		return false, err
	}

	deviceReq, err := http.NewRequest("PUT", "https://api.spotify.com/v1/me/player", bytes.NewReader(b))
	if err != nil {
		return false, err
	}

	deviceReq.Header.Set("Authorization", "Bearer "+accessToken)
	deviceReq.Header.Set("Content-Type", "application/json")
	deviceRes, err := s.httpClient.Do(deviceReq)
	if err != nil {
		return false, err
	}

	switch deviceRes.StatusCode {
	case http.StatusNoContent:
		return true, nil
	case http.StatusBadRequest:
		return false, errors.New(pifyErrors.BAD_OR_EXPIRED_TOKEN)
	case http.StatusForbidden:
		return false, errors.New(pifyErrors.BAD_OAUTH_REQUEST)
	case http.StatusTooManyRequests:
		return false, errors.New(pifyErrors.RATE_LIMIT_EXCEEDED)
	default:
		return false, errors.New(pifyErrors.UNKNOWN_ERROR)
	}
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
