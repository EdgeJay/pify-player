package models

import (
	"time"

	"github.com/uptrace/bun"
)

// TrackMedia represents additional media for song track.
type TrackMedia struct {
	bun.BaseModel

	Id             int64 `bun:",pk,autoincrement"`
	SpotifyTrackId string
	MediaType      string
	MediaId        string
	CreatedAt      time.Time  `bun:",notnull,default:current_timestamp"`
	DeletedAt      *time.Time `bun:",soft_delete"`
}
