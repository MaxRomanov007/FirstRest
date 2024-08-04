package car

import (
	"context"
	"errors"
	"fmt"
	"server/internal/domain/models"
	"server/internal/storage/postgres/cars"
)

func (s *Service) FindCar(
	ctx context.Context,
	number string,
) (models.Car, error) {
	const op = "services.cars.car.Find"

	number = formatLicensePlate(number)

	car, err := s.prov.Find(ctx, number)
	if errors.Is(err, cars.ErrNotFound) {
		return models.Car{}, fmt.Errorf("%s: %w", op, ErrCarNotFound)
	}
	if err != nil {
		return models.Car{}, fmt.Errorf("%s: %w", op, err)
	}

	return car, nil
}
