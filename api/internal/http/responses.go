package http

type UserDetails struct {
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
	IsController    *bool  `json:"is_controller"`
}

type LoginResponse struct {
	LoggedIn    bool         `json:"logged_in"`
	User        *UserDetails `json:"user"`
	RedirectUrl string       `json:"redirect_url"`
	ErrorCode   string       `json:"error_code"`
}

type YoutubeVideoResponse struct {
	VideoId string `json:"video_id"`
}

type ApiResponse struct {
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"error_code"`
}
