package images

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func (s *Service) UpdateImages(
	images []*multipart.FileHeader,
	number string,
) error {
	const op = "services.cars.images.UpdateImages"

	number = formatLicensePlate(number)

	if len(images) == 0 {
		return fmt.Errorf("%s: %w", op, ErrNoOneImage)
	}

	if len(images) > 15 {
		return fmt.Errorf("%s: %w", op, ErrTooManyImages)
	}

	dir := filepath.Join(s.cfg.ImagesPath, number)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%s: %s: %w", op, dir, ErrNotFound)
	}

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%s: %s: fail remove directory: %w", op, dir, err)
	}

	err = os.MkdirAll(dir, os.ModeDir)
	if os.IsExist(err) {
		return fmt.Errorf("%s: %s: dicectory still exists: %w", op, dir, err)
	}
	if err != nil {
		return fmt.Errorf("%s: %s: fail to create directory: %w", op, dir, err)
	}

	for i, _ := range images {
		file, err := images[i].Open()
		if err != nil {
			return fmt.Errorf("%s: %s: fail to open image: %w", op, dir, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Printf("%s: failed to close file: %s\n", op, err.Error())
			}
		}()

		fileName := strconv.Itoa(i) + filepath.Ext(images[i].Filename)
		filePath := filepath.Join(dir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("%s: %s: fail to create destination file: %w", op, filePath, err)
		}
		defer func() {
			err := dst.Close()
			if err != nil {
				fmt.Printf("%s: failed to close file: %s\n", op, err.Error())
			}
		}()

		if _, err := io.Copy(dst, file); err != nil {
			return fmt.Errorf("%s: %s: fail to copy file: %w", op, fileName, err)
		}
	}
	return nil
}
