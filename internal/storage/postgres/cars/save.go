package cars

import (
	"context"
	"fmt"
	"github.com/lib/pq"
)

func (s *Storage) Save(
	ctx context.Context,
	producer string,
	model string,
	engineCapacity float32,
	power float32,
	number string,
	imagesCount uint8,
	description string,
) (int64, error) {
	const op = "postgres.cars.Save"

	stmt, err := s.db.Prepare(`
INSERT INTO cars 
    (producer, 
     model, 
     engine_capacity, 
     power, 
     number, 
     images_count, 
     description) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
RETURNING id
`)
	if err != nil {
		return 0, fmt.Errorf("%s: failed prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(
		ctx,
		producer,
		model,
		engineCapacity,
		power,
		number,
		imagesCount,
		description,
	)

	var id int64

	err = row.Scan(&id)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
		return 0, fmt.Errorf("%s: %w", op, ErrCarExists)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed scan row: %w", op, err)
	}

	return id, nil
}
