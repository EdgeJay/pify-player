package errors

// auth related error codes
const (
	MISSING_CODE_OR_STATE      = "missing_code_or_state"
	GET_ACCESS_TOKEN_FAILED    = "get_access_token_failed"
	GET_USER_INFO_FAILED       = "get_user_info_failed"
	SAVE_USER_INFO_FAILED      = "save_user_info_failed"
	GENERATE_SESSION_ID_FAILED = "generate_session_id_failed"
	SAVE_SESSION_FAILED        = "save_session_failed"
)

// api related error codes
const (
	UNKNOWN_ERROR               = "unknown_error"
	GET_DEVICES_FAILED          = "get_devices_failed"
	INVALID_REQUEST_BODY        = "invalid_request_body"
	INVALID_SESSION             = "invalid_session"
	LOGIN_QR_UNAVAILABLE        = "login_qr_unavailable"
	BAD_OR_EXPIRED_TOKEN        = "bad_or_expired_token"
	BAD_OAUTH_REQUEST           = "bad_oauth_request"
	RATE_LIMIT_EXCEEDED         = "rate_limit_exceeded"
	GET_TRACK_FAILED            = "get_track_failed"
	PARSE_TRACK_RESPONSE_FAILED = "parse_track_response_failed"
	NO_YOUTUBE_VIDEO_FOUND      = "no_youtube_video_found"
)
