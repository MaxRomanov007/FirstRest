package car

import (
	"context"
	"fmt"
	"server/internal/domain/models"
	"strconv"
	"strings"
)

func (s *Service) GetCars(
	ctx context.Context,
	filter models.Filter,
) ([]models.Car, string, error) {
	const op = "services.cars.car.Get"

	if filter.OrderBy == "" {
		filter.OrderBy = s.cfg.DefaultOrderBy
	}

	cars, totalCount, err := s.prov.Get(
		ctx,
		strings.ToLower(filter.SearchQuery),
		filter.Limit,
		filter.Page,
		filter.OrderBy,
		filter.Desc,
	)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return cars, strconv.FormatInt(totalCount, strconvBase), nil
}
