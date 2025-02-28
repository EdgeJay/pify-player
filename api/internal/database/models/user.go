package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	Id        int64          `bun:",pk,autoincrement"`
	Username  string         `bun:"type:,unique"`
	Sessions  []*UserSession `bun:"rel:has-many,join:id=user_id"`
	CreatedAt time.Time      `bun:",notnull,default:current_timestamp"`
	DeletedAt *time.Time     `bun:",soft_delete"`
}
