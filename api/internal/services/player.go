package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/database/models"
	"github.com/uptrace/bun"

	"github.com/dgraph-io/ristretto/v2"
)

const (
	PLAYER_STATE_DISCONNECTED = "disconnected"
	PLAYER_STATE_WAITING      = "waiting"
	PLAYER_STATE_CONNECTED    = "connected"
)

type PlayerService struct {
	db    *database.SQLiteDB
	cache *ristretto.Cache[string, string]
}

func NewPlayerService(db *database.SQLiteDB) *PlayerService {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	ps := &PlayerService{db, cache}
	ps.initialize()

	return ps
}

func (s *PlayerService) initialize() {
	s.cache.Set("player_state", PLAYER_STATE_DISCONNECTED, 1)
	s.cache.Wait()
}

func (s *PlayerService) GetPlayerState() string {
	state, found := s.cache.Get("player_state")
	if !found {
		panic(errors.New("missing player state in cache"))
	}
	return state
}

func (s *PlayerService) Connect() (*models.UserSession, error) {
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
		if err == sql.ErrNoRows {
			s.cache.Set("player_state", PLAYER_STATE_WAITING, 1)
			s.cache.Wait()
			return nil, nil
		}
		return nil, err
	}

	s.cache.Set("player_state", PLAYER_STATE_CONNECTED, 1)
	s.cache.Wait()

	return session, nil
}
