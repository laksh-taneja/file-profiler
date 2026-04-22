package utils

import (
	"errors"
	"os"
)

func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	// Permission denied etc.
	return false, err
}
