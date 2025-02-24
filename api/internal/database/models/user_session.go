package models

import (
	"time"

	"github.com/uptrace/bun"
)

type UserSession struct {
	bun.BaseModel

	Id                    int64  `bun:",pk,autoincrement"`
	UserId                int64  `bun:",notnull"`
	Uuid                  string `bun:",notnull"`
	UserAgent             string `bun:",notnull"`
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	CreatedAt             time.Time `bun:",notnull,default:current_timestamp"`
	DeletedAt             time.Time `bun:",soft_delete"`
}
