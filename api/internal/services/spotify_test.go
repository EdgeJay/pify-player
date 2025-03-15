package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Helper type to rewrite URLs for testing
type rewriteTransport struct {
	URL       string
	NewURL    string
	Transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface. It intercepts HTTP requests
// and rewrites the URL if it matches the configured URL pattern, then delegates
// the actual HTTP transport to the underlying Transport. If the request URL matches
// t.URL, it will be rewritten to t.NewURL before being forwarded.
func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.String() == t.URL {
		req.URL, _ = url.Parse(t.NewURL)
	}
	return t.Transport.RoundTrip(req)
}

func TestGetAuthUrl(t *testing.T) {
	// Setup test service
	credentials := SpotifyCredentials{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}
	service := NewSpotifyService(credentials, &http.Client{})

	// Call GetAuthUrl
	authUrl, err := service.GetAuthUrl()
	if err != nil {
		t.Fatalf("GetAuthUrl returned error: %v", err)
	}

	// Parse the returned URL
	parsed, err := url.Parse(authUrl)
	if err != nil {
		t.Fatalf("Failed to parse auth URL: %v", err)
	}

	// Verify query parameters
	query := parsed.Query()

	if query.Get("client_id") != credentials.ClientID {
		t.Errorf("Expected client_id %s, got %s", credentials.ClientID, query.Get("client_id"))
	}

	if query.Get("response_type") != "code" {
		t.Errorf("Expected response_type 'code', got %s", query.Get("response_type"))
	}

	if query.Get("redirect_uri") != credentials.RedirectURI {
		t.Errorf("Expected redirect_uri %s, got %s", credentials.RedirectURI, query.Get("redirect_uri"))
	}

	if query.Get("state") == "" {
		t.Error("State parameter is empty")
	}

	if len(query.Get("state")) != 16 {
		t.Errorf("Expected state length 16, got %d", len(query.Get("state")))
	}

	if query.Get("scope") == "" {
		t.Error("Scope parameter is empty")
	}
}

func TestGetApiToken(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify content type header
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		// Verify auth header
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("No basic auth provided")
		}
		if username != "test-client-id" || password != "test-client-secret" {
			t.Errorf("Wrong basic auth credentials: %s:%s", username, password)
		}

		// Verify form data
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "authorization_code" {
			t.Errorf("Expected grant_type authorization_code, got %s", r.Form.Get("grant_type"))
		}
		if r.Form.Get("code") != "test-code" {
			t.Errorf("Expected code test-code, got %s", r.Form.Get("code"))
		}
		if r.Form.Get("redirect_uri") != "http://localhost:8080/callback" {
			t.Errorf("Expected redirect_uri http://localhost:8080/callback, got %s", r.Form.Get("redirect_uri"))
		}

		// Send response
		response := SpotifyTokenResponse{
			AccessToken:  "test-access-token",
			ExpiresIn:    3600,
			RefreshToken: "test-refresh-token",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Setup test service
	credentials := SpotifyCredentials{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}
	client := mockServer.Client()
	service := NewSpotifyService(credentials, client)

	// Override token URL to use mock server
	originalTokenURL := "https://accounts.spotify.com/api/token"
	client.Transport = rewriteTransport{
		URL:       originalTokenURL,
		NewURL:    mockServer.URL,
		Transport: http.DefaultTransport,
	}

	// Test GetApiToken
	token, err := service.GetApiToken("test-code")
	if err != nil {
		t.Fatalf("GetApiToken returned error: %v", err)
	}

	// Verify response
	if token.AccessToken != "test-access-token" {
		t.Errorf("Expected access token test-access-token, got %s", token.AccessToken)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("Expected expires_in 3600, got %d", token.ExpiresIn)
	}
	if token.RefreshToken != "test-refresh-token" {
		t.Errorf("Expected refresh token test-refresh-token, got %s", token.RefreshToken)
	}
}

func TestRefreshApiToken(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify content type header
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		// Verify form data
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Errorf("Expected grant_type refresh_token, got %s", r.Form.Get("grant_type"))
		}
		if r.Form.Get("refresh_token") != "test-refresh-token" {
			t.Errorf("Expected refresh_token test-refresh-token, got %s", r.Form.Get("refresh_token"))
		}
		if r.Form.Get("client_id") != "test-client-id" {
			t.Errorf("Expected client_id test-client-id, got %s", r.Form.Get("client_id"))
		}

		// Send response
		response := SpotifyTokenResponse{
			AccessToken:  "new-access-token",
			ExpiresIn:    3600,
			RefreshToken: "new-refresh-token",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Setup test service
	credentials := SpotifyCredentials{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}
	client := mockServer.Client()
	service := NewSpotifyService(credentials, client)

	// Override token URL to use mock server
	originalTokenURL := "https://accounts.spotify.com/api/token"
	client.Transport = rewriteTransport{
		URL:       originalTokenURL,
		NewURL:    mockServer.URL,
		Transport: http.DefaultTransport,
	}

	// Test RefreshApiToken
	token, err := service.RefreshApiToken("test-refresh-token")
	if err != nil {
		t.Fatalf("RefreshApiToken returned error: %v", err)
	}

	// Verify response
	if token.AccessToken != "new-access-token" {
		t.Errorf("Expected access token new-access-token, got %s", token.AccessToken)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("Expected expires_in 3600, got %d", token.ExpiresIn)
	}
	if token.RefreshToken != "new-refresh-token" {
		t.Errorf("Expected refresh token new-refresh-token, got %s", token.RefreshToken)
	}
}

func TestGetUser(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify authorization header
		if r.Header.Get("Authorization") != "Bearer test-access-token" {
			t.Errorf("Expected Authorization Bearer test-access-token, got %s", r.Header.Get("Authorization"))
		}

		// Send response
		response := SpotifyUser{
			Id:          "test-user-id",
			DisplayName: "Test User",
			Images: []struct {
				Url    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			}{
				{
					Url:    "https://example.com/image.jpg",
					Width:  300,
					Height: 300,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Setup test service
	credentials := SpotifyCredentials{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}
	client := mockServer.Client()
	service := NewSpotifyService(credentials, client)

	// Override API URL to use mock server
	originalUserURL := "https://api.spotify.com/v1/me"
	client.Transport = rewriteTransport{
		URL:       originalUserURL,
		NewURL:    mockServer.URL,
		Transport: http.DefaultTransport,
	}

	// Test GetUser
	user, err := service.GetUser("test-access-token")
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}

	// Verify response
	if user.Id != "test-user-id" {
		t.Errorf("Expected user ID test-user-id, got %s", user.Id)
	}
	if user.DisplayName != "Test User" {
		t.Errorf("Expected display name Test User, got %s", user.DisplayName)
	}
	if len(user.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(user.Images))
	}
	if user.Images[0].Url != "https://example.com/image.jpg" {
		t.Errorf("Expected image URL https://example.com/image.jpg, got %s", user.Images[0].Url)
	}
	if user.Images[0].Width != 300 {
		t.Errorf("Expected image width 300, got %d", user.Images[0].Width)
	}
	if user.Images[0].Height != 300 {
		t.Errorf("Expected image height 300, got %d", user.Images[0].Height)
	}
}
