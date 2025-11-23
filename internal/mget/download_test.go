package mget

import (
	"testing"
)

func TestDownload_InvalidURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Invalid URL format",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Unsupported platform",
			url:     "https://example.com/video",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Download(tt.url, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownload_ValidYouTubeURL(t *testing.T) {
	// Note: This test will fail if there's no internet connection or if YouTube is unavailable
	// In a real scenario, you'd want to mock the YouTube client
	// For now, we'll test that the URL validation works correctly

	validURLs := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ",
		"https://m.youtube.com/watch?v=dQw4w9WgXcQ",
	}

	for _, url := range validURLs {
		t.Run("Valid URL: "+url, func(t *testing.T) {
			// We can't actually download without mocking, but we can verify URL validation
			isValid, platform := validatePlatformURL(url)
			if !isValid {
				t.Errorf("Expected URL %q to be valid", url)
			}
			if platform != "youtube" {
				t.Errorf("Expected platform 'youtube', got %q", platform)
			}
		})
	}
}

func TestDownload_ErrorMessages(t *testing.T) {
	// Test invalid URL error message
	err := Download("not-a-url", false)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
	if err.Error() != "invalid video URL" {
		t.Errorf("Expected 'invalid video URL' error, got %q", err.Error())
	}

	// Test unsupported platform error message
	err = Download("https://example.com/video", false)
	if err == nil {
		t.Error("Expected error for unsupported platform")
	}
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}
