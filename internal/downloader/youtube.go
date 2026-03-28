package downloader

import (
	"fmt"
	"regexp"
)

type YouTubeProvider struct {}

func (p *YouTubeProvider) CanHandle(url string) bool {
	// Valid formats:
	// - Standard/Mobile: https://(www.|m.)youtube.com/watch?v=ID
	// - Shortened: https://youtu.be/ID
	// - Shorts: https://youtube.com/shorts/ID
	// - Embed/Legacy: https://youtube.com/(embed|v)/ID
	var youtubeRegex = regexp.MustCompile(`(?i)^(?:https?://)?(?:www\.|m\.)?(?:youtu\.be/|youtube\.com/(?:watch\?v=|embed/|shorts/|v/))([a-zA-Z0-9_-]{11})`)
	return youtubeRegex.MatchString(url)
}

func (p *YouTubeProvider) Handle(url string) error {
  fmt.Printf("downloading youtube video from %s\n", url)
  return nil
}
