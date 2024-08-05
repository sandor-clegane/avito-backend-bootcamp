package postgres

import (
	"avito-backend-bootcamp/internal/config"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func buildDSN(cfg *config.DB) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", cfg.Username, cfg.Password, cfg.Name, cfg.Host, cfg.Port)
}

func New(ctx context.Context, cfg *config.DB) (*Repository, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", buildDSN(cfg))
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, migrationQuery)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}
