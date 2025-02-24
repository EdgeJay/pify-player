package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/edgejay/pify-player/api/internal/database/models"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().
			Model((*models.UserSession)(nil)).
			IfNotExists().
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().
			Model((*models.UserSession)(nil)).
			Exec(ctx)
		return err
	})
}
