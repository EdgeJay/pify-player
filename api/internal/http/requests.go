package http

type YoutubeVideoRequest struct {
	Query          string `json:"query"`
	SpotifyTrackId string `json:"spotify_track_id"`
	CacheResults   bool   `json:"cache_results"`
}

type PlayerCommandRequest struct {
	Command string `json:"command"`
}
