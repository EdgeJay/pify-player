package utils

import (
	"os"
	"strings"
)

type ServerSettings struct {
	Port        string
	CorsOrigins []string
	SslDomain   string
}

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func GetServerSettings() ServerSettings {
	return ServerSettings{
		Port:        os.Getenv("PORT"),
		CorsOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		SslDomain:   os.Getenv("SSL_DOMAIN"),
	}
}

func GetDBFilename() string {
	return os.Getenv("DB_FILE")
}

func GetSpotifyCredentials() SpotifyCredentials {
	return SpotifyCredentials{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
	}
}

func GetCallbackDestination() string {
	return os.Getenv("CALLBACK_DEST")
}
