package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestGetServerSettings(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected ServerSettings
	}{
		{
			name: "basic settings",
			envVars: map[string]string{
				"PORT":         "8080",
				"CORS_ORIGINS": "http://localhost:3000,http://example.com",
				"SSL_DOMAIN":   "example.com",
			},
			expected: ServerSettings{
				Port:        "8080",
				CorsOrigins: []string{"http://localhost:3000", "http://example.com"},
				SslDomain:   "example.com",
			},
		},
		{
			name: "empty settings",
			envVars: map[string]string{
				"PORT":         "",
				"CORS_ORIGINS": "",
				"SSL_DOMAIN":   "",
			},
			expected: ServerSettings{
				Port:        "",
				CorsOrigins: []string{""},
				SslDomain:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			got := GetServerSettings()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetServerSettings() = %v, want %v", got, tt.expected)
			}
		})
	}
}
