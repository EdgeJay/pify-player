package services

import (
	"context"
	"database/sql"

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
