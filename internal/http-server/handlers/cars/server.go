package cars

import (
	"context"
	"log/slog"
	"mime/multipart"
	"server/internal/config"
	"server/internal/domain/models"
)

type CarService interface {
	SaveCar(
		ctx context.Context,
		car models.Car,
	) (id int64, err error)

	DeleteCar(
		ctx context.Context,
		number string,
	) (id int64, err error)

	UpdateCar(
		ctx context.Context,
		car models.Car,
	) (id int64, err error)

	FindCar(
		ctx context.Context,
		number string,
	) (car models.Car, err error)

	GetCars(
		ctx context.Context,
		filter models.Filter, //orderBy: "producer" | "model" | "engine_capacity" | "power" | "images_count"
	) (cars []models.Car, totalCount string, err error)
}

type ImagesService interface {
	GetImage(
		number string,
		imgId string,
	) (path string, totalCount int8, err error)

	SaveImages(
		images []*multipart.FileHeader,
		number string,
	) (err error)

	UpdateImages(
		images []*multipart.FileHeader,
		number string,
	) (err error)

	DeleteImages(
		number string,
	) (err error)
}

type API struct {
	Log   *slog.Logger
	Cfg   *config.CarsConfig
	Car   CarService
	Image ImagesService
}

const (
	strconvBase               = 10
	strconvBitSize            = 64
	totalCountHeader          = "X-Total-Count"
	accessControlExposeHeader = "Access-Control-Expose-Headers"
)

func New(
	log *slog.Logger,
	cfg *config.CarsConfig,
	carService CarService,
	ImgService ImagesService,
) *API {
	return &API{
		Log:   log,
		Cfg:   cfg,
		Car:   carService,
		Image: ImgService,
	}
}
