package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/database/models"
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
		Scan(context.Background())

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SaveUser(username string) (*models.User, error) {
	// check if user by username exists
	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Return existing user if found
	if user != nil && user.Id > 0 {
		return user, nil
	}

	// or else create new user
	_, err = s.db.Bun.NewInsert().
		Model(&models.User{Username: username, DeletedAt: nil}).
		Exec(context.Background())

	if err != nil {
		return nil, err
	}

	// Get user from database
	user, err = s.GetUser(username)
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
		Exists(context.Background())
}

func (s *UserService) GetSession(sessionId string) (*models.UserSession, error) {
	// fetch session
	session := &models.UserSession{}

	// Get user from database
	if err := s.db.Bun.NewSelect().
		Model(session).
		Where("uuid = ?", sessionId).
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
