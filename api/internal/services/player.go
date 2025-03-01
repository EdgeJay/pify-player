package services

import (
	"context"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/database/models"
)

type PlayerService struct {
	db *database.SQLiteDB
}

func NewPlayerService(db *database.SQLiteDB) *PlayerService {
	return &PlayerService{db}
}

func (s *PlayerService) ClearPreviousStates() error {
	_, err := s.db.Bun.NewDelete().
		Model((*models.PlayerState)(nil)).
		Where("is_active = ?", true).
		WhereOr("deleted_at IS NULL").
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (s *PlayerService) SetWaitingState() error {
	// clear previous statuses first
	err := s.ClearPreviousStates()
	if err != nil {
		return err
	}

	_, err = s.db.Bun.NewInsert().
		Model(&models.PlayerState{IsWaiting: true, IsActive: false, DeletedAt: nil}).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
