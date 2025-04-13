package http

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectResponse(t *testing.T) {
	response := ConnectResponse{
		LoginResponse: LoginResponse{
			LoggedIn:    true,
			User:        &UserDetails{DisplayName: "Test User", ProfileImageUrl: "http://example.com/image.jpg"},
			RedirectUrl: "http://example.com/redirect",
			ErrorCode:   "",
		},
		Connected: true,
	}

	if response.LoggedIn != true {
		t.Errorf("Expected LoggedIn to be true, got %v", response.LoggedIn)
	}
	if response.User.DisplayName != "Test User" {
		t.Errorf("Expected DisplayName to be 'Test User', got %s", response.User.DisplayName)
	}
	if response.Connected != true {
		t.Errorf("Expected Connected to be true, got %v", response.Connected)
	}

	b, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Error marshaling response: %v", err)
	}

	assert.Contains(t, string(b), `"logged_in":true`)
	assert.Contains(t, string(b), `"display_name":"Test User"`)
	assert.Contains(t, string(b), `"connected":true`)
	assert.Contains(t, string(b), `"redirect_url":"http://example.com/redirect"`)
	assert.Contains(t, string(b), `"error_code":""`)
	assert.Contains(t, string(b), `"profile_image_url":"http://example.com/image.jpg"`)
}
