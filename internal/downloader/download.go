package downloader

import (
	"fmt"

	"github.com/renato3x/mget/internal/utils"
)

func Download(url string) error {
  if !utils.IsValidUrl(url) {
    return fmt.Errorf("Invalid url: %s", url)
  }

  return nil
}
