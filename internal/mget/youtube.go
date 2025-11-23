package mget

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
)

type mediaDownloadParams struct {
	video      *youtube.Video
	client     youtube.Client
	outputDir  string
	mimeType   string
	extension  string
}

type progressDownloadParams struct {
	stream              io.ReadCloser
	file                *os.File
	contentLength       int64
	formatContentLength int64
	videoTitle          string
	filepath            string
}

func DownloadYoutube(url string, audioOnly bool) error {
	// Fetch youtube video
	client := youtube.Client{}
	video, err := client.GetVideo(url)
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

func downloadWithProgress(params progressDownloadParams) error {
	// Calculate total bytes: use contentLength from stream, fallback to format ContentLength if available
	totalBytes := params.contentLength
	if totalBytes == 0 && params.formatContentLength > 0 {
		totalBytes = params.formatContentLength
	}

	bar := ProgressWriter(totalBytes, "Downloading media from YouTube...")

	if _, err := io.Copy(io.MultiWriter(params.file, bar), params.stream); err != nil {
		return err
	}

	// Print success message
	fmt.Println()
	fmt.Printf("Media \"%s\" downloaded to %s\n", params.videoTitle, params.filepath)

	return nil
}

func downloadVideo(video *youtube.Video, client youtube.Client, outputDir string) error {
	params := mediaDownloadParams{
		video:     video,
		client:    client,
		outputDir: outputDir,
		mimeType:  "video/mp4",
		extension: ".mp4",
	}
	return downloadMedia(params)
}

func downloadAudio(video *youtube.Video, client youtube.Client, outputDir string) error {
	params := mediaDownloadParams{
		video:     video,
		client:    client,
		outputDir: outputDir,
		mimeType:  "audio/mp4",
		extension: ".mp3",
	}
	return downloadMedia(params)
}

func downloadMedia(params mediaDownloadParams) error {
	// Get formats with audio channels
	formats := params.video.Formats.WithAudioChannels()
	var selectedFormat *youtube.Format

	// Find format matching the requested MIME type
	for _, format := range formats {
		if strings.Contains(format.MimeType, params.mimeType) {
			selectedFormat = &format
			break
		}
	}

	if selectedFormat == nil {
		return fmt.Errorf("no valid format found")
	}

	// Get stream
	stream, contentLength, err := params.client.GetStream(params.video, selectedFormat)
	if err != nil {
		return fmt.Errorf("error downloading media")
	}
	defer stream.Close()

	// Create file
	filename := GenerateFilename()
	filepath := filepath.Join(params.outputDir, filename+params.extension)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error downloading media")
	}
	defer file.Close()

	// Download with progress
	progressParams := progressDownloadParams{
		stream:              stream,
		file:                file,
		contentLength:       contentLength,
		formatContentLength: selectedFormat.ContentLength,
		videoTitle:          params.video.Title,
		filepath:            filepath,
	}
	if err := downloadWithProgress(progressParams); err != nil {
		return fmt.Errorf("error downloading media")
	}

	return nil
}
