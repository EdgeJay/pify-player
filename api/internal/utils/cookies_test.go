package utils

import (
	"net/http"
	"testing"
	"time"
)

func TestCreateCookie(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		duration time.Time
	}{
		{
			name:     "test_cookie",
			value:    "test_value",
			duration: time.Now().Add(24 * time.Hour),
		},
		{
			name:     "empty_cookie",
			value:    "",
			duration: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := CreateCookie(tt.name, tt.value, tt.duration)

			if cookie.Name != tt.name {
				t.Errorf("CreateCookie().Name = %v, want %v", cookie.Name, tt.name)
			}
			if cookie.Value != tt.value {
				t.Errorf("CreateCookie().Value = %v, want %v", cookie.Value, tt.value)
			}
			if !cookie.Expires.Equal(tt.duration) {
				t.Errorf("CreateCookie().Expires = %v, want %v", cookie.Expires, tt.duration)
			}
			if cookie.Path != "/" {
				t.Errorf("CreateCookie().Path = %v, want /", cookie.Path)
			}
			if cookie.SameSite != http.SameSiteNoneMode {
				t.Errorf("CreateCookie().SameSite = %v, want %v", cookie.SameSite, http.SameSiteNoneMode)
			}
			if !cookie.HttpOnly {
				t.Error("CreateCookie().HttpOnly = false, want true")
			}
			if !cookie.Secure {
				t.Error("CreateCookie().Secure = false, want true")
			}
		})
	}
}
