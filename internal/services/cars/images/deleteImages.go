package images

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *Service) DeleteImages(number string) error {
	const op = "services.cars.images.DeleteImages"

	number = formatLicensePlate(number)

	dir := filepath.Join(s.cfg.ImagesPath, number)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%s: %s: %w", op, dir, ErrNotFound)
	}

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%s: %s: fail remove directory: %w", op, dir, err)
	}

	return nil
}
