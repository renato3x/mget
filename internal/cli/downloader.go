package cli

import (
	"fmt"
	"net/url"
)

func Download(url, output string, audioOnly bool) error {
	isValid, urlData := isValidURL(url)

	if !isValid {
		return fmt.Errorf("invalid video URL")
	}

	if !isValidWebsite(urlData.Host) {
		return fmt.Errorf("unsupported website. accepted sites: %s", getAcceptedSites())
	}

	return nil
}

var allowedHosts = map[string]struct{}{
	"youtube.com":    {},
	"youtu.be":       {},
	"tiktok.com":     {},
	"www.tiktok.com": {},
	"vm.tiktok.com":  {},
}

func isValidWebsite(host string) bool {
	_, ok := allowedHosts[host]
	return ok
}

func isValidURL(str string) (bool, *url.URL) {
	data, err := url.ParseRequestURI(str)
	return err == nil, data
}

func getAcceptedSites() string {
	return "YouTube, TikTok"
}
