package images

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *Service) GetImage(
	number string,
	imgId string,
) (string, int8, error) {
	const op = "services.cars.images.GetImage"

	number = formatLicensePlate(number)

	dir := filepath.Join(s.cfg.ImagesPath, number)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", 0, fmt.Errorf("%s: %s: %w", op, dir, ErrNotFound)
	}

	var imgPaths []string
	var imgCount int8

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		imgCount++

		if fileNameWithoutExtension(info.Name()) == imgId {
			imgPaths = append(imgPaths, path)
		}

		return nil
	})

	if err != nil {
		return "", 0, fmt.Errorf("%s: %s: failed to walk directory: %w", op, dir, err)
	}

	if len(imgPaths) == 0 {
		return "", 0, fmt.Errorf("%s: %s: %w", op, dir, ErrNotFound)
	}

	return imgPaths[0], imgCount, nil
}
