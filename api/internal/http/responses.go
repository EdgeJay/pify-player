package responses

type UserDetails struct {
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type LoginResponse struct {
	LoggedIn    bool         `json:"logged_in"`
	User        *UserDetails `json:"user"`
	RedirectUrl string       `json:"redirect_url"`
	ErrorCode   string       `json:"error_code"`
}
