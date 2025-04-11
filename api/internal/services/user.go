package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/uptrace/bun"
)

type UserService struct {
	db *database.SQLiteDB
}

func NewUserService(db *database.SQLiteDB) *UserService {
	return &UserService{db}
}

func (s *UserService) GetUser(username string) (*models.User, error) {
	user := &models.User{}

	// Get user from database
	err := s.db.Bun.NewSelect().
		Model(user).
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(context.Background())

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SaveUser(spotifyUser *SpotifyUser) (*models.User, error) {
	// check if user by username exists
	user, err := s.GetUser(spotifyUser.Id)
	if err != nil {
		return nil, err
	}

	profileImageUrl := ""
	if len(spotifyUser.Images) > 0 {
		profileImageUrl = spotifyUser.Images[0].Url
	}

	// Update existing user record if found
	if user != nil && user.Id > 0 {
		_, err := s.db.Bun.NewUpdate().
			Model((*models.User)(nil)).
			Set("display_name = ?", spotifyUser.DisplayName).
			Set("profile_image_url = ?", profileImageUrl).
			Where("username = ?", spotifyUser.Id).
			Exec(context.Background())

		if err != nil {
			return nil, err
		}
	} else {
		// or else create new user
		_, err := s.db.Bun.NewInsert().
			Model(&models.User{
				Username:        spotifyUser.Id,
				DisplayName:     spotifyUser.DisplayName,
				ProfileImageUrl: profileImageUrl,
				DeletedAt:       nil,
			}).
			Exec(context.Background())

		if err != nil {
			return nil, err
		}
	}

	// Get user from database
	user, err = s.GetUser(spotifyUser.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SessionExists(sessionId string) (bool, error) {
	// check if session exists
	session := &models.UserSession{}

	return s.db.Bun.NewSelect().
		Model(session).
		Where("uuid = ?", sessionId).
		Where("deleted_at IS NULL").
		Exists(context.Background())
}

func (s *UserService) GetSession(sessionId string) (*models.UserSession, error) {
	// fetch session
	session := &models.UserSession{}

	// Get user from database
	if err := s.db.Bun.NewSelect().
		Model(session).
		Relation("User", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("display_name", "profile_image_url")
		}).
		Where("user_session.uuid = ?", sessionId).
		Where("user_session.deleted_at IS NULL").
		Scan(context.Background()); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *UserService) SaveSession(
	userId int64,
	sessionId,
	userAgent,
	accessToken,
	refreshToken string,
	accessTokenExpiresAt time.Time,
) (*models.UserSession, error) {
	// check if session exists
	exists, err := s.SessionExists(sessionId)

	if err != nil {
		return nil, err
	}

	if !exists {
		// create new session
		_, err = s.db.Bun.NewInsert().
			Model(&models.UserSession{
				UserId:                userId,
				Uuid:                  sessionId,
				UserAgent:             userAgent,
				AccessToken:           accessToken,
				RefreshToken:          refreshToken,
				AccessTokenExpiresAt:  accessTokenExpiresAt,
				RefreshTokenExpiresAt: nil,
				DeletedAt:             nil,
			}).
			Exec(context.Background())

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	// fetch session
	return s.GetSession(sessionId)
}

func (s *UserService) UpdateSessionAccessToken(
	sessionId,
	accessToken string,
	accessTokenExpiresAt time.Time,
) (*models.UserSession, error) {
	_, err := s.db.Bun.NewUpdate().
		Model((*models.UserSession)(nil)).
		Set("access_token = ?", accessToken).
		Set("access_token_expires_at = ?", accessTokenExpiresAt).
		Where("uuid = ?", sessionId).
		Exec(context.Background())

	if err != nil {
		return nil, err
	}

	// fetch session
	return s.GetSession(sessionId)
}

func (s *UserService) DeleteSession(sessionId string) error {
	_, err := s.db.Bun.NewDelete().
		Model((*models.UserSession)(nil)).
		Where("uuid = ?", sessionId).
		Exec(context.Background())
	return err
}
