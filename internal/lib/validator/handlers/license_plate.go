package handlers

import (
	"github.com/go-playground/validator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

func LicensePlateValidator(fl validator.FieldLevel) bool {
	return LicensePlateValidatorByString(fl.Field().String())
}

func LicensePlateValidatorByString(str string) bool {
	str = cases.Upper(language.Und).String(str)

	match, _ := regexp.MatchString(
		`^[ABEKMHOPCTYXАВЕКМНОРСТУХ][0-9]{3}[ABEKMHOPCTYXАВЕКМНОРСТУХ]{2}$`,
		str,
	)
	if !match {
		return false
	}

	if strings.Contains(str, "000") {
		return false
	}

	return true
}

func New() *validator.Validate {
	val := validator.New()
	if err := val.RegisterValidation("license_plate", LicensePlateValidator); err != nil {
		return validator.New()
	}
	return val
}
