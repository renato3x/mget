package mget

import (
	"os"
	"path/filepath"

  "github.com/google/uuid"
)

func GetOutputDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "mget-downloads"), nil
}

func GenerateFilename() string {
  return uuid.New().String()
}
