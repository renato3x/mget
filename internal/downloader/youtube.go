package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	survey "github.com/AlecAivazis/survey/v2"
	youtube "github.com/kkdai/youtube/v2"
	"github.com/renato3x/mget/internal/utils"
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
  client := youtube.Client{}
  video, err := client.GetVideo(url)

  if err != nil {
    return fmt.Errorf("Some error occurred searching media")
  }

  fmt.Printf("%s - %s\n", video.Title, video.Duration)

  optionsWithAudio := video.Formats.WithAudioChannels()
  var options []youtube.Format
  for len(options) == 0 {
    options, err = p.selectOptions(optionsWithAudio)
    if err != nil {
      return err
    }
  }

  return p.downloadOptions(options, &client, video)
}

func (p *YouTubeProvider) selectOptions(formats youtube.FormatList) ([]youtube.Format, error) {
  var options []string
  formatMap := make(map[string]*youtube.Format)

  for _, f := range formats {
    label := ""
    if f.Height > 0 {
      label = fmt.Sprintf("[%d] Audio + Video - %s (%s)", f.ItagNo, f.QualityLabel, f.MimeType)
    } else {
      label = fmt.Sprintf("[%d] Audio only - (%s)", f.ItagNo, f.MimeType)
    }
    options = append(options, label)
    formatMap[label] = &f
  }

  selectedOptions := []string{}
  prompt := &survey.MultiSelect{
    Message: "Select one or more options for download",
    Options: options,
  }

  err := survey.AskOne(prompt, &selectedOptions)

  if err != nil {
    return nil, fmt.Errorf("Error selecting download options")
  }

  selectedFormats := []youtube.Format{}
  for _, selection := range selectedOptions {
    selectedFormats = append(selectedFormats, *formatMap[selection])
  }

  return selectedFormats, nil
}

func (p *YouTubeProvider) downloadOptions(
  formats []youtube.Format,
  client *youtube.Client,
  video *youtube.Video,
) error {
  outputDir, err := GetOutputDirectory()

  if err != nil {
    return fmt.Errorf("Error getting output directory for downloads")
  }

  if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("Error creating output directory")
	}

  for _, f := range formats {
    err := func(format youtube.Format) error {
      fmt.Printf("Downloading: %d - %s...\n", format.ItagNo, format.MimeType)
      stream, _, err := client.GetStream(video, &format)

      if err != nil {
        return err
      }
      defer stream.Close()

      filename := utils.Slugify(fmt.Sprintf("%d %s %s", format.ItagNo, format.MimeType, video.Title))
      ext := utils.GetFileExtensionByMimetype(format.MimeType)
      filepath := filepath.Join(outputDir, filename+ext)
      file, err := os.Create(filepath)

      if err != nil {
        return err
      }
      defer file.Close()

      _, err = io.Copy(file, stream)
      if err != nil {
        return err
      }

      fmt.Printf("%s (%s) downloaded successfully\n", video.Title, format.MimeType)
      
      return nil
    }(f)
    
    if err != nil {
      fmt.Printf("Some error occurred downloading option %d - %s\n", f.ItagNo, f.MimeType)
    }
  }

  return nil
}
