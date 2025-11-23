package mget

import "fmt"

func Download(url string, audioOnly bool) error {
	isValid, platform := validatePlatformURL(url)

	if !isValid {
		return fmt.Errorf("invalid video URL")
	}

	if platform == "" {
		return fmt.Errorf("unsupported website. accepted sites: %s", getAcceptedSites())
	}

	return nil
}
