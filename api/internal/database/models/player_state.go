package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PlayerState struct {
	bun.BaseModel

	Id        int64 `bun:",pk,autoincrement"`
	UserId    int64
	IsWaiting bool       `bun:",default:false"`
	IsActive  bool       `bun:",default:false"`
	CreatedAt time.Time  `bun:",notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:",soft_delete"`
}
