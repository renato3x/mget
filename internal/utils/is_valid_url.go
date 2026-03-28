package utils

import "net/url"

func IsValidUrl(rawUrl string) bool {
	_, err := url.ParseRequestURI(rawUrl)
	return err == nil
}
