package car

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/postgres/cars"
)

func (s *Service) DeleteCar(
	ctx context.Context,
	number string,
) (int64, error) {
	const op = "services.cars.car.DeleteCar"

	number = formatLicensePlate(number)

	id, err := s.own.Delete(ctx, number)
	if errors.Is(err, cars.ErrNotFound) {
		return 0, fmt.Errorf("%s: %w", op, ErrCarNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
