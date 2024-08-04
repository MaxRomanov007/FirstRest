package response

import (
	"fmt"
	"github.com/go-playground/validator"
	"net/http"
	"strings"
)

func InternalError(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func AlreadyExists(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusConflict)
}

func NF(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func ValidationError(w http.ResponseWriter, errs validator.ValidationErrors) {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is required", err.Field()))
		case "license_plate":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is not valid license plate", err.Field()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is over maximum", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is not valid", err.Field()))
		}
	}

	http.Error(w, strings.Join(errMsgs, "; "), http.StatusBadRequest)
}
