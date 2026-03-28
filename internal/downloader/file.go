package downloader

import (
	"os"
	"path/filepath"
)

func GetOutputDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "mget-downloads"), nil
}
