package mget

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func DownloadYoutube(url string, audioOnly bool) error {
	// Extract video ID from URL
	videoID, err := extractYoutubeVideoID(url)
	if err != nil {
		return err
	}

	// Fetch youtube video
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		return fmt.Errorf("error fetching youtube video")
	}

	// Get output directory
	outputDir, err := GetOutputDirectory()
	if err != nil {
		return fmt.Errorf("error getting output directory")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory")
	}

	// Route to appropriate download function
	if audioOnly {
		return downloadAudio(video, client, outputDir)
	}

	return downloadVideo(video, client, outputDir)
}

func downloadVideo(video *youtube.Video, client youtube.Client, outputDir string) error {
	// Get the best audio format
	formats := video.Formats.WithAudioChannels()
	format := &formats[0]

  for _, frmt := range formats {
    fmt.Println(frmt.MimeType)
  }

	stream, contentLength, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("error downloading youtube video")
	}

	defer stream.Close()

	// Create file
	filename := GenerateFilename()
	filepath := filepath.Join(outputDir, filename+".mp4")

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error downloading youtube video")
	}

	defer file.Close()

	// Download with progress
	// Use contentLength from stream, fallback to format ContentLength if available
	totalBytes := contentLength
	if totalBytes == 0 && format.ContentLength > 0 {
		totalBytes = format.ContentLength
	}

	bar := ProgressWriter(totalBytes, "Downloading")

	if _, err = io.Copy(io.MultiWriter(file, bar), stream); err != nil {
		return fmt.Errorf("error downloading youtube video")
	}

	// Print success message
  fmt.Println()
	fmt.Printf("Video \"%s\" downloaded to %s\n", video.Title, filepath)

	return nil
}

func downloadAudio(video *youtube.Video, client youtube.Client, outputDir string) error {
  formats := video.Formats.WithAudioChannels()
  var selectedFormat *youtube.Format

  for _, format := range formats {
    if strings.Contains(format.MimeType, "audio/mp4") {
      selectedFormat = &format
      break
    }
  }

  if selectedFormat == nil {
    return fmt.Errorf("no audio format found")
  }

  stream, contentLength, err := client.GetStream(video, selectedFormat)
  if err != nil {
    return fmt.Errorf("error downloading youtube audio")
  }

  defer stream.Close()

  // Create file
  filename := GenerateFilename()
  filepath := filepath.Join(outputDir, filename+".mp3")

  file, err := os.Create(filepath)
  if err != nil {
    return fmt.Errorf("error downloading youtube audio")
  }

  defer file.Close()

  // Download with progress
	// Use contentLength from stream, fallback to format ContentLength if available
	totalBytes := contentLength
	if totalBytes == 0 && selectedFormat.ContentLength > 0 {
		totalBytes = selectedFormat.ContentLength
	}

	bar := ProgressWriter(totalBytes, "Downloading")

	if _, err = io.Copy(io.MultiWriter(file, bar), stream); err != nil {
		return fmt.Errorf("error downloading youtube audio")
	}

	// Print success message
  fmt.Println()
	fmt.Printf("Video \"%s\" downloaded to %s\n", video.Title, filepath)

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
