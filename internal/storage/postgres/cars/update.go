package cars

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) Update(
	ctx context.Context,
	producer string,
	model string,
	engineCapacity float32,
	power float32,
	number string,
	imagesCount uint8,
	description string,
) (int64, error) {
	const op = "postgres.cars.Update"

	stmt, err := s.db.Prepare(`
UPDATE cars 
SET producer = $1, 
    model = $2, 
    engine_capacity = $3, 
    power = $4, 
    images_count = $5, 
    description = $6 
WHERE number = $7 
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
		imagesCount,
		description,
		number,
	)

	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed scan row: %w", op, err)
	}

	return id, nil
}
