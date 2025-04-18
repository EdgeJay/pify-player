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

func GetServerSettings() ServerSettings {
	return ServerSettings{
		Port:        os.Getenv("SERVER_PORT"),
		CorsOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		SslDomain:   os.Getenv("SSL_DOMAIN"),
	}
}

func GetDBFilename() string {
	return os.Getenv("DB_FILE")
}

func GetCallbackDestination() string {
	return os.Getenv("CALLBACK_DEST")
}

func GetBasicAuthUsername() string {
	return os.Getenv("BASIC_AUTH_USERNAME")
}

func GetBasicAuthPassword() string {
	return os.Getenv("BASIC_AUTH_PASSWORD")
}

func GetYoutubeApiKey() string {
	return os.Getenv("YOUTUBE_API_KEY")
}

func ShellCommandsAllowed() bool {
	return os.Getenv("ALLOW_SHELL_COMMANDS") == "1"
}
