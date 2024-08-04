package main

import (
	"errors"
	"flag"
	"github.com/fatih/color"
	"log"
	"os"
	"server/internal/config"
	pgApi "server/internal/lib/api/storage/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	color.NoColor = false

	var (
		cfgPath string
		migPath string
		vers    uint
	)

	flag.StringVar(&cfgPath, "config_path", "", "path to config file")
	flag.StringVar(&migPath, "migrations_path", "", "path to config file")
	flag.UintVar(&vers, "version", 0, "migration version")
	flag.Parse()

	cfg := config.MustLoadByPath(cfgPath)

	m, err := migrate.New(
		"file://"+migPath,
		pgApi.ConnString(&cfg.DB),
	)
	if err != nil {
		log.Fatal("failed to create migrator: " + err.Error())
	}

	if vers == 0 {
		if err := m.Up(); err != nil {
			notNilMigrateErr(err)
		}
	} else {
		if err := m.Migrate(vers); err != nil {
			notNilMigrateErr(err)
		}
	}

	color.Green("\nmigrations applied successfully\n\n")
}

func notNilMigrateErr(err error) {
	if errors.Is(err, migrate.ErrNoChange) {
		color.Blue("\nno migrations to apply\n\n")

		os.Exit(2)
	}

	log.Fatal("failed to confirm migrations: " + err.Error())
}
