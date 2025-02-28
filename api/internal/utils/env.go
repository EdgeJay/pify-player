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
		Port:        os.Getenv("PORT"),
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
