package responses

type LoginResponse struct {
	LoggedIn    bool   `json:"logged_in"`
	RedirectUrl string `json:"redirect_url"`
	ErrorCode   string `json:"error_code"`
}
