package postgres

import (
	repo "avito-backend-bootcamp/internal/infra/repository"
	"database/sql"

	"github.com/lib/pq"
)

// https://github.com/lib/pq/blob/master/error.go
func PostgresErrorTransform(err error) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(*pq.Error)
	if ok {
		if pgErr.Code == "23505" {
			return repo.ErrAlreadyExists
		}
		if pgErr.Code == "23503" || pgErr.Code == "23502" {
			return repo.ErrConstraintViolation
		}
	}

	if err == sql.ErrNoRows {
		return repo.ErrNotFound
	}

	return err
}
