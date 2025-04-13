package services

import (
	"context"
	"errors"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/database/models"
	pifyErrors "github.com/edgejay/pify-player/api/internal/errors"
	"github.com/uptrace/bun"
)

const (
	PLAYER_STATE_DISCONNECTED = "disconnected"
	PLAYER_STATE_WAITING      = "waiting"
	PLAYER_STATE_CONNECTED    = "connected"
)

type TrackMediaType string

const (
	TRACK_MEDIA_TYPE_YOUTUBE TrackMediaType = "youtube"
)

type PlayerService struct {
	db *database.SQLiteDB
}

func NewPlayerService(db *database.SQLiteDB) *PlayerService {
	ps := &PlayerService{db}
	return ps
}

func (s *PlayerService) GetControllerSession() (*models.UserSession, error) {
	// fetch session
	session := &models.UserSession{}
	err := s.db.Bun.NewSelect().
		Model(session).
		Relation("User", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("display_name", "profile_image_url")
		}).
		Where("user_session.is_controller = TRUE").
		Scan(context.Background())

	if err != nil {
		// no session or session with controller privileges found
		return nil, errors.New(pifyErrors.INVALID_SESSION)
	}

	return session, nil
}

func (s *PlayerService) GetTrackMedia(spotifyTrackId string, mediaType TrackMediaType) *models.TrackMedia {
	trackMedia := &models.TrackMedia{}

	err := s.db.Bun.NewSelect().
		Model(trackMedia).
		Where("spotify_track_id = ? AND media_type = ?", spotifyTrackId, mediaType).
		Scan(context.Background())

	if err != nil {
		return nil
	}

	return trackMedia
}

func (s *PlayerService) SaveTrackMedia(spotifyTrackId, mediaId string, mediaType TrackMediaType) error {
	_, err := s.db.Bun.NewInsert().
		Model(&models.TrackMedia{
			SpotifyTrackId: spotifyTrackId,
			MediaId:        mediaId,
			MediaType:      string(mediaType),
		}).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
