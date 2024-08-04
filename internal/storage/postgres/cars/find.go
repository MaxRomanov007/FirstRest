package cars

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"server/internal/domain/models"
)

func (s *Storage) Find(
	ctx context.Context,
	number string,
) (models.Car, error) {
	const op = "postgres.cars.Find"

	stmt, err := s.db.Prepare(`
SELECT 
    producer, 
    model, 
    engine_capacity, 
    power, 
    number, 
    images_count, 
    description 
FROM cars
WHERE number = $1
`)
	if err != nil {
		return models.Car{}, fmt.Errorf("%s: failed prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, number)

	var car models.Car
	err = row.Scan(
		&car.Producer,
		&car.Model,
		&car.EngineCapacity,
		&car.Power,
		&car.Number,
		&car.ImagesCount,
		&car.Description,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Car{}, ErrNotFound
	}
	if err != nil {
		return models.Car{}, fmt.Errorf("%s: failed scan row: %w", op, err)
	}

	return car, nil
}
