package cars

import (
	"context"
	"fmt"
	"server/internal/domain/models"
)

func (s *Storage) Get(
	ctx context.Context,
	searchQuery string,
	limit int64,
	page int64,
	orderBy string, //"producer", "model", "engine_capacity", "power", "images_count"
	desc bool,
) ([]models.Car, int64, error) {
	const op = "postgres.cars.Get"

	//Getting total count...

	stmt, err := s.db.Prepare(`SELECT COUNT(*) FROM cars`)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: failed prepare get total count statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx)

	var totalCount int64
	err = row.Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: failed get total count row: %w", op, err)
	}

	//Getting cars...

	var descStr string
	if desc {
		descStr = "desc"
	} else {
		descStr = "asc"
	}

	stmt, err = s.db.Prepare(`
SELECT 
    producer,
    model,
    engine_capacity,
    power,
    number,
    images_count,
    description
FROM cars
WHERE Lower(producer) || ' ' || Lower(model) || ' ' || Lower(description) || Lower(number) LIKE '%' || $1 || '%' 
ORDER BY ` + orderBy + " " + descStr + ` 
limit $2
offset $3
`)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: failed prepare get cars statement: %w", op, err)
	}

	if page != 0 {
		page = page - 1
	}
	page *= limit

	rows, err := stmt.QueryContext(ctx, searchQuery, limit, page)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: failed using get cars statement: %w", op, err)
	}

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(
			&car.Producer,
			&car.Model,
			&car.EngineCapacity,
			&car.Power,
			&car.Number,
			&car.ImagesCount,
			&car.Description,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: failed scan car row: %w", op, err)
		}
		cars = append(cars, car)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return cars, totalCount, nil
}
