package car

import (
	"context"
	"errors"
	"fmt"
	"server/internal/domain/models"
	"server/internal/storage/postgres/cars"
)

func (s *Service) UpdateCar(
	ctx context.Context,
	car models.Car,
) (int64, error) {
	const op = "services.cars.car.UpdateCar"

	car.Number = formatLicensePlate(car.Number)

	id, err := s.own.Update(
		ctx,
		car.Producer,
		car.Model,
		car.EngineCapacity,
		car.Power,
		car.Number,
		car.ImagesCount,
		car.Description,
	)
	if errors.Is(err, cars.ErrNotFound) {
		return 0, ErrCarExists
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
