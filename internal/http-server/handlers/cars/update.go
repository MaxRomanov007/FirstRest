package cars

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/domain/models"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	customValidator "server/internal/lib/validator/handlers"
	"server/internal/services/cars/car"
	images2 "server/internal/services/cars/images"
	"strings"
)

type UpdateRequest struct {
	Producer       string  `validate:"required,max=255" json:"producer"`
	Model          string  `validate:"required,max=255" json:"model"`
	EngineCapacity float32 `validate:"required,number" json:"engine_capacity"`
	Power          float32 `validate:"required,number" json:"power"`
	Number         string  `validate:"required,license_plate" json:"number"`
	Description    string  `validate:"omitempty,required" json:"description,omitempty"`
}

func (c *API) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cars.Update"

		var reqId = middleware.GetReqID(r.Context())

		log := c.Log.With(
			slog.String("operation", op),
			slog.String("request_id", reqId),
		)

		err := r.ParseMultipartForm(10000)
		if err != nil {
			log.Error("failed to parse multipart form", sl.Err(err))

			response.InternalError(w)

			return
		}

		m := r.MultipartForm

		var req UpdateRequest

		reader := strings.NewReader(m.Value["body"][0])
		if err := render.DecodeJSON(reader, &req); err != nil {
			log.Error("failed to decode request body")

			response.InternalError(w)

			return
		}

		if err := customValidator.New().Struct(&req); err != nil {
			var validError validator.ValidationErrors
			if ok := errors.As(err, &validError); ok {
				log.Warn("invalid request", sl.Err(err))

				response.ValidationError(w, validError)

				return
			}

			log.Error("failed to valid request", sl.Err(err))

			response.InternalError(w)

			return
		}

		images := m.File[c.Cfg.ImagesFormName]
		err = c.Image.UpdateImages(images, req.Number)
		if errors.Is(err, images2.ErrNotFound) {
			log.Warn("car not found", sl.Err(err))

			response.NF(w, r)

			return
		}
		if errors.Is(err, images2.ErrTooManyImages) {
			log.Warn("too many images", sl.Err(err))

			response.BadRequest(w, "too many images")

			return
		}
		if errors.Is(err, images2.ErrNoOneImage) {
			log.Warn("no one image", sl.Err(err))

			response.BadRequest(w, "no one image")

			return
		}
		if err != nil {
			log.Error("failed to save images", sl.Err(err))

			response.InternalError(w)

			return
		}

		log.Debug("updating car...")

		_, err = c.Car.UpdateCar(
			r.Context(),
			models.Car{
				Producer:       req.Producer,
				Model:          req.Model,
				EngineCapacity: req.EngineCapacity,
				Power:          req.Power,
				Number:         req.Number,
				ImagesCount:    uint8(len(images)),
				Description:    req.Description,
			})
		if errors.Is(err, car.ErrCarNotFound) {
			log.Warn("car not found", sl.Err(err))

			response.NF(w, r)

			return
		}
		if err != nil {
			log.Error("failed to save cars", sl.Err(err))

			response.InternalError(w)

			return
		}
	}
}
