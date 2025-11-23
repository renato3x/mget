package mget

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func DownloadYoutube(url string, audioOnly bool) error {
	_, err := extractYoutubeVideoID(url)

	if err != nil {
		return err
	}

	return nil
}

func extractYoutubeVideoID(youtubeURL string) (string, error) {
	parsedURL, err := url.Parse(youtubeURL)
	if err != nil {
		return "", fmt.Errorf("invalid YouTube URL")
	}

	host := strings.ToLower(parsedURL.Host)
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")

	// For URLs like youtu.be/VIDEO_ID
	if host == "youtu.be" {
		videoID := strings.TrimPrefix(parsedURL.Path, "/")
		// Remove query parameters if present
		if idx := strings.Index(videoID, "?"); idx != -1 {
			videoID = videoID[:idx]
		}
		if idx := strings.Index(videoID, "&"); idx != -1 {
			videoID = videoID[:idx]
		}
		if isValidVideoID(videoID) {
			return videoID, nil
		}
	}

	// For URLs like youtube.com/watch?v=VIDEO_ID
	if host == "youtube.com" || host == "m.youtube.com" {
		// Try to get from query parameter 'v'
		if videoID := parsedURL.Query().Get("v"); videoID != "" && isValidVideoID(videoID) {
			return videoID, nil
		}

		// For URLs like youtube.com/embed/VIDEO_ID
		if strings.HasPrefix(parsedURL.Path, "/embed/") {
			videoID := strings.TrimPrefix(parsedURL.Path, "/embed/")
			videoID = strings.TrimPrefix(videoID, "/")
			if idx := strings.Index(videoID, "?"); idx != -1 {
				videoID = videoID[:idx]
			}
			if isValidVideoID(videoID) {
				return videoID, nil
			}
		}

		// For URLs like youtube.com/v/VIDEO_ID
		if strings.HasPrefix(parsedURL.Path, "/v/") {
			videoID := strings.TrimPrefix(parsedURL.Path, "/v/")
			videoID = strings.TrimPrefix(videoID, "/")
			if idx := strings.Index(videoID, "?"); idx != -1 {
				videoID = videoID[:idx]
			}
			if isValidVideoID(videoID) {
				return videoID, nil
			}
		}
	}

	return "", fmt.Errorf("could not extract video ID from URL")
}

func isValidVideoID(videoID string) bool {
	if len(videoID) != 11 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", videoID)
	return matched
}
