package cars

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) Delete(
	ctx context.Context,
	number string,
) (int64, error) {
	const op = "postgres.cars.Delete"

	stmt, err := s.db.Prepare(`
DELETE FROM cars
WHERE number = $1
RETURNING id
`)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, number)

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
