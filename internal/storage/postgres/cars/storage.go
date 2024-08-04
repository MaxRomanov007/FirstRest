package cars

import (
	"database/sql"
	"errors"
	"fmt"
	"server/internal/config"

	pgApi "server/internal/lib/api/storage/postgres"
)

var (
	ErrCarExists = errors.New("car already exists")
	ErrNotFound  = errors.New("car not found")
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.DbConnConfig) (*Storage, error) {
	const op = "postgres.cars.New"

	db, err := sql.Open("postgres", pgApi.ConnString(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}
