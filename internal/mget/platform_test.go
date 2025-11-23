package mget

import (
	"testing"
)

func TestValidatePlatformURL(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		expectedValid    bool
		expectedPlatform string
	}{
		{
			name:             "Valid YouTube URL - youtube.com",
			url:              "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expectedValid:    true,
			expectedPlatform: "youtube",
		},
		{
			name:             "Valid YouTube URL - youtu.be",
			url:              "https://youtu.be/dQw4w9WgXcQ",
			expectedValid:    true,
			expectedPlatform: "youtube",
		},
		{
			name:             "Valid YouTube URL - m.youtube.com",
			url:              "https://m.youtube.com/watch?v=dQw4w9WgXcQ",
			expectedValid:    true,
			expectedPlatform: "youtube",
		},
		{
			name:             "Valid YouTube URL - www.youtube.com",
			url:              "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expectedValid:    true,
			expectedPlatform: "youtube",
		},
		{
			name:             "Valid YouTube URL with port",
			url:              "https://youtube.com:443/watch?v=dQw4w9WgXcQ",
			expectedValid:    true,
			expectedPlatform: "youtube",
		},
		{
			name:             "Invalid URL - not a URL",
			url:              "not-a-url",
			expectedValid:    false,
			expectedPlatform: "",
		},
		{
			name:             "Valid URL but unsupported platform",
			url:              "https://example.com/video",
			expectedValid:    true,
			expectedPlatform: "",
		},
		{
			name:             "Empty URL",
			url:              "",
			expectedValid:    false,
			expectedPlatform: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, platform := validatePlatformURL(tt.url)

			if valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectedValid, valid)
			}
			if platform != tt.expectedPlatform {
				t.Errorf("Expected platform=%q, got %q", tt.expectedPlatform, platform)
			}
		})
	}
}

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple host",
			input:    "youtube.com",
			expected: "youtube.com",
		},
		{
			name:     "Host with www",
			input:    "www.youtube.com",
			expected: "youtube.com",
		},
		{
			name:     "Host with port",
			input:    "youtube.com:443",
			expected: "youtube.com",
		},
		{
			name:     "Host with www and port",
			input:    "www.youtube.com:443",
			expected: "youtube.com",
		},
		{
			name:     "Uppercase host",
			input:    "YOUTUBE.COM",
			expected: "youtube.com",
		},
		{
			name:     "Mixed case host",
			input:    "YouTube.Com",
			expected: "youtube.com",
		},
		{
			name:     "Subdomain",
			input:    "m.youtube.com",
			expected: "m.youtube.com",
		},
		{
			name:     "Subdomain with www",
			input:    "www.m.youtube.com",
			expected: "m.youtube.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeHost(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIdentifyPlatform(t *testing.T) {
	tests := []struct {
		name             string
		host             string
		expectedPlatform string
	}{
		{
			name:             "YouTube host - youtube.com",
			host:             "youtube.com",
			expectedPlatform: "youtube",
		},
		{
			name:             "YouTube host - youtu.be",
			host:             "youtu.be",
			expectedPlatform: "youtube",
		},
		{
			name:             "YouTube host - m.youtube.com",
			host:             "m.youtube.com",
			expectedPlatform: "youtube",
		},
		{
			name:             "YouTube subdomain",
			host:             "subdomain.youtube.com",
			expectedPlatform: "youtube",
		},
		{
			name:             "Unsupported host",
			host:             "example.com",
			expectedPlatform: "",
		},
		{
			name:             "Empty host",
			host:             "",
			expectedPlatform: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := identifyPlatform(tt.host)
			if result != tt.expectedPlatform {
				t.Errorf("Expected platform %q, got %q", tt.expectedPlatform, result)
			}
		})
	}
}

func TestGetAcceptedSites(t *testing.T) {
	result := getAcceptedSites()

	// Should contain "youtube" at minimum
	if result == "" {
		t.Error("Expected non-empty accepted sites string")
	}

	// Should contain youtube
	if len(result) < len("youtube") {
		t.Error("Expected accepted sites to contain platform names")
	}
}
