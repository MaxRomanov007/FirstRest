package postgres

import (
	"fmt"
	"server/internal/config"
)

func ConnString(cfg *config.DbConnConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Server, cfg.Port, cfg.DB,
	)
}
