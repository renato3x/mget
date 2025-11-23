package mget

import (
	"net/url"
	"strings"
)

var platformHosts = map[string][]string{
	"youtube": {
		"youtube.com",
		"youtu.be",
		"m.youtube.com",
	},
}

func validatePlatformURL(str string) (bool, string) {
	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		return false, ""
	}

	host := normalizeHost(parsedURL.Host)

	platform := identifyPlatform(host)
	if platform != "" {
		return true, platform
	}

	return true, ""
}

func normalizeHost(str string) string {
  host := strings.ToLower(str)
	
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Remove www. to normalize
	host = strings.TrimPrefix(host, "www.")

  return host
}

func identifyPlatform(host string) string {
	for platform, hosts := range platformHosts {
		for _, h := range hosts {
			if host == h || strings.HasSuffix(host, "."+h) {
				return platform
			}
		}
	}

	return ""
}

func getAcceptedSites() string {
	platforms := make([]string, 0, len(platformHosts))
	for platform := range platformHosts {
		platforms = append(platforms, platform)
	}

	return strings.Join(platforms, ", ")
}
