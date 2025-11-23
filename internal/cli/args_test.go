package cli

import (
	"flag"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedUrl     string
		expectedAudio   bool
		expectedVersion bool
	}{
		{
			name:            "URL only",
			args:            []string{"https://youtube.com/watch?v=test"},
			expectedUrl:     "https://youtube.com/watch?v=test",
			expectedAudio:   false,
			expectedVersion: false,
		},
		{
			name:            "URL with audio flag",
			args:            []string{"-a", "https://youtube.com/watch?v=test"},
			expectedUrl:     "https://youtube.com/watch?v=test",
			expectedAudio:   true,
			expectedVersion: false,
		},
		{
			name:            "URL with --audio flag",
			args:            []string{"--audio", "https://youtube.com/watch?v=test"},
			expectedUrl:     "https://youtube.com/watch?v=test",
			expectedAudio:   true,
			expectedVersion: false,
		},
		{
			name:            "Version flag only",
			args:            []string{"-v"},
			expectedUrl:     "",
			expectedAudio:   false,
			expectedVersion: true,
		},
		{
			name:            "Version flag with --version",
			args:            []string{"--version"},
			expectedUrl:     "",
			expectedAudio:   false,
			expectedVersion: true,
		},
		{
			name:            "No arguments",
			args:            []string{},
			expectedUrl:     "",
			expectedAudio:   false,
			expectedVersion: false,
		},
		{
			name:            "Audio and version flags",
			args:            []string{"-a", "-v", "https://youtube.com/watch?v=test"},
			expectedUrl:     "https://youtube.com/watch?v=test",
			expectedAudio:   true,
			expectedVersion: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Reset flag package state
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Set test args
			os.Args = append([]string{"mget"}, tt.args...)

			// Parse args
			args := Args()

			// Verify results
			if args.Url != tt.expectedUrl {
				t.Errorf("Expected URL %q, got %q", tt.expectedUrl, args.Url)
			}
			if args.AudioOnly != tt.expectedAudio {
				t.Errorf("Expected AudioOnly %v, got %v", tt.expectedAudio, args.AudioOnly)
			}
			if args.Version != tt.expectedVersion {
				t.Errorf("Expected Version %v, got %v", tt.expectedVersion, args.Version)
			}
		})
	}
}
