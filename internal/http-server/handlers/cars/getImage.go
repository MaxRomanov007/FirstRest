package cars

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	customValidator "server/internal/lib/validator/handlers"
	"server/internal/services/cars/images"
	"strconv"
)

type getImageRequest struct {
	Number string `json:"number" validate:"required,license_plate"`
	Id     string `json:"id" validate:"required,number,max=2"`
}

func (c *API) GetImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cars.GetImage"

		log := c.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req getImageRequest

		req.Number = r.URL.Query().Get("number")
		req.Id = r.URL.Query().Get("id")

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

		imgPath, totalCount, err := c.Image.GetImage(req.Number, req.Id)
		if errors.Is(err, images.ErrNotFound) {
			log.Warn("image not found", sl.Err(err))

			response.NF(w, r)

			return
		}
		if err != nil {
			log.Error("failed getting image", sl.Err(err))

			response.InternalError(w)

			return
		}

		w.Header().Add(
			accessControlExposeHeader,
			totalCountHeader,
		)
		if w.Header().Get(totalCountHeader) == "" {
			w.Header().Add(
				totalCountHeader,
				strconv.Itoa(int(totalCount)),
			)
		}

		http.ServeFile(w, r, imgPath)
	}
}
