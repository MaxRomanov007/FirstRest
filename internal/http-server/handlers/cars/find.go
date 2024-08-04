package cars

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	customValidator "server/internal/lib/validator/handlers"
	carService "server/internal/services/cars/car"
)

type findRequest struct {
	Number string `json:"number" validate:"required,license_plate"`
}

func (c *API) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cars.find"

		var reqId = middleware.GetReqID(r.Context())

		log := c.Log.With(
			slog.String("operation", op),
			slog.String("request_id", reqId),
		)

		req := findRequest{
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

		log.Debug("finding car...")

		car, err := c.Car.FindCar(r.Context(), req.Number)
		if errors.Is(err, carService.ErrCarNotFound) {
			log.Warn("car not found")

			response.NotFound(w, "car not found")

			return
		}

		log.Info("car found")

		render.JSON(w, r, car)
	}
}
