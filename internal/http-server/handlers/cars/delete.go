package cars

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	customValidator "server/internal/lib/validator/handlers"
	"server/internal/services/cars/car"
)

type deleteRequest struct {
	Number string `json:"number" validate:"required,license_plate"`
}

func (c *API) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cars.Delete"

		var reqId = middleware.GetReqID(r.Context())

		log := c.Log.With(
			slog.String("operation", op),
			slog.String("request_id", reqId),
		)

		number := chi.URLParam(r, "number")
		if number == "" {
			log.Warn("number is empty")

			response.BadRequest(w, "number is empty")

			return
		}

		req := deleteRequest{
			Number: chi.URLParam(r, "number"),
		}

		if err := customValidator.New().Struct(req); err != nil {
			var validErr validator.ValidationErrors
			if errors.As(err, &validErr) {
				log.Warn("invalid number", sl.Err(err))

				response.ValidationError(w, validErr)

				return
			}

			log.Error("failed to validate request", sl.Err(err))

			response.InternalError(w)

			return
		}

		log.Debug("deleting car...")

		id, err := c.Car.DeleteCar(r.Context(), number)
		if errors.Is(err, car.ErrCarNotFound) {
			log.Warn("car not found", sl.Err(err))

			response.NF(w, r)

			return
		}
		if err != nil {
			log.Error("failed to delete car", sl.Err(err))

			response.InternalError(w)

			return
		}

		log.Info("car deleted", slog.Int64("id", id))

		err = c.Image.DeleteImages(number)
		if err != nil {
			log.Error("failed to delete images", sl.Err(err))

			response.InternalError(w)

			return
		}
	}
}
