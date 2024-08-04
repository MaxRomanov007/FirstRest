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
	"strconv"
)

type getRequest struct {
	limit       string `validate:"omitempty,number"`
	page        string `validate:"omitempty,number"`
	orderBy     string `validate:"omitempty,oneof=producer model engine_capacity power images_count"`
	desc        string `validate:"omitempty,boolean"`
	searchQuery string `validate:"omitempty"`
}

type parsedGetRequest struct {
	limit       int64
	page        int64
	orderBy     string
	desc        bool
	searchQuery string
}

func (c *API) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cars.Get"

		log := c.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req getRequest

		req.limit = r.URL.Query().Get("limit")
		req.page = r.URL.Query().Get("page")
		req.orderBy = r.URL.Query().Get("order_by")
		req.desc = r.URL.Query().Get("desc")
		req.searchQuery = r.URL.Query().Get("query")

		if err := validator.New().Struct(req); err != nil {
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
		if req.limit == "" {
			req.limit = "0"
		}
		if req.page == "" {
			req.page = "0"
		}
		if req.desc == "" {
			req.desc = "false"
		}

		var parsedReq parsedGetRequest
		var err error

		parsedReq.orderBy = req.orderBy
		parsedReq.searchQuery = req.searchQuery
		parsedReq.limit, err = strconv.ParseInt(req.limit, strconvBase, strconvBitSize)
		if err != nil {
			log.Warn("failed to convert 'limit' to int64", sl.Err(err))

			response.InternalError(w)

			return
		}
		parsedReq.page, err = strconv.ParseInt(req.page, strconvBase, strconvBitSize)
		if err != nil {
			log.Warn("failed to convert 'page' to int64", sl.Err(err))

			response.InternalError(w)

			return
		}
		parsedReq.desc, err = strconv.ParseBool(req.desc)
		if err != nil {
			log.Warn("failed to convert 'desc' to bool", sl.Err(err))

			response.InternalError(w)

			return
		}

		log.Debug("getting cars...")

		cars, totalCount, err := c.Car.GetCars(
			r.Context(),
			models.Filter{
				SearchQuery: parsedReq.searchQuery,
				Limit:       parsedReq.limit,
				Page:        parsedReq.page,
				OrderBy:     parsedReq.orderBy,
				Desc:        parsedReq.desc,
			},
		)
		if err != nil {
			log.Error("failed to get cars", sl.Err(err))

			response.InternalError(w)

			return
		}

		log.Info("cars got")

		w.Header().Add(
			accessControlExposeHeader,
			totalCountHeader,
		)
		if w.Header().Get(totalCountHeader) == "" {
			w.Header().Add(
				totalCountHeader,
				totalCount,
			)
		}

		render.JSON(w, r, cars)
	}
}
