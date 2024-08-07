package postgres

import (
	"avito-backend-bootcamp/internal/config"
	dbUtil "avito-backend-bootcamp/pkg/utils/db"
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func New(ctx context.Context, cfg *config.DB) (*Repository, error) {
	db, err := sqlx.ConnectContext(
		ctx, "postgres",
		dbUtil.BuildDSN(cfg.Username, cfg.Password, cfg.Name, cfg.Host, cfg.Port),
	)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:     db,
		getter: trmsqlx.DefaultCtxGetter,
	}, nil
}
