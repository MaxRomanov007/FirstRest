package car

import (
	"context"
	"errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"server/internal/config"
	"server/internal/domain/models"
	"strings"
)

// Owner describes storage functions create, update and delete
type Owner interface {
	Save(
		ctx context.Context,
		producer string,
		model string,
		engineCapacity float32,
		power float32,
		number string,
		imagesCount uint8,
		description string,
	) (id int64, err error)

	Update(
		ctx context.Context,
		producer string,
		model string,
		engineCapacity float32,
		power float32,
		number string,
		imagesCount uint8,
		description string,
	) (id int64, err error)

	Delete(
		ctx context.Context,
		number string,
	) (id int64, err error)
}

type Provider interface {
	Find(
		ctx context.Context,
		number string,
	) (car models.Car, err error)

	Get(
		ctx context.Context,
		searchQuery string,
		limit int64,
		page int64,
		orderBy string, //"producer", "model", "engine_capacity", "power", "images_count"
		desc bool,
	) (cars []models.Car, totalCount int64, err error)
}

const (
	strconvBase = 10
)

var (
	ErrCarExists   = errors.New("cars already exists")
	ErrCarNotFound = errors.New("cars not found")
)

type Service struct {
	cfg  *config.CarsConfig
	own  Owner
	prov Provider
}

func New(cfg *config.CarsConfig, own Owner, prov Provider) *Service {
	return &Service{
		cfg:  cfg,
		own:  own,
		prov: prov,
	}
}

/*
formatLicensePlate replaces russian characters to
english in license plate and make it in upper case
*/
func formatLicensePlate(old string) string {
	replacer := strings.NewReplacer(
		"А", "A",
		"В", "B",
		"Е", "E",
		"К", "K",
		"М", "M",
		"Н", "H",
		"О", "O",
		"Р", "P",
		"С", "C",
		"Т", "T",
		"У", "Y",
		"Х", "X",
	)

	strToReplace := cases.Upper(language.Und).String(old)

	return replacer.Replace(strToReplace)
}
