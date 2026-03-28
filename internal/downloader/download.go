package downloader

import (
	"fmt"

	"github.com/renato3x/mget/internal/utils"
)

var providers = []Provider{
	&YouTubeProvider{},
}

func Download(url string) error {
	if !utils.IsValidUrl(url) {
		return fmt.Errorf("Invalid url: %s", url)
	}

	for _, provider := range providers {
		if provider.CanHandle(url) {
			return provider.Handle(url)
		}
	}

	return fmt.Errorf("No provider available for url: %s", url)
}
