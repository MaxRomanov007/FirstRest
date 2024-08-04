package images

import (
	"errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"server/internal/config"
	"strings"
)

var (
	ErrTooManyImages = errors.New("too many images")
	ErrNoOneImage    = errors.New("no one image")
	ErrNotFound      = errors.New("image not found")
	ErrExists        = errors.New("images directory already exists")
)

type Service struct {
	cfg *config.CarsConfig
}

func New(cfg *config.CarsConfig) *Service {
	return &Service{cfg: cfg}
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
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
