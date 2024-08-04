package app

import (
	"fmt"
	"log/slog"
	carsapp "server/internal/app/cars"
	"server/internal/config"
	carServer "server/internal/http-server/handlers/cars"
	carService "server/internal/services/cars/car"
	"server/internal/services/cars/images"
	carStorage "server/internal/storage/postgres/cars"
)

type App struct {
	CarsServer *carsapp.App
}

func MustLoad(log *slog.Logger, cfg *config.Config) *App {
	app, err := New(log, cfg)
	if err != nil {
		panic(err)
	}

	return app
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	const op = "app.New"

	carsStorage, err := carStorage.New(&cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	carsCarsService := carService.New(&cfg.Cars, carsStorage, carsStorage)
	carsImagesService := images.New(&cfg.Cars)

	carsServer := carServer.New(log, &cfg.Cars, carsCarsService, carsImagesService)

	carsApp := carsapp.New(&cfg.HTTPServer, carsServer)

	return &App{
		CarsServer: carsApp,
	}, nil
}
