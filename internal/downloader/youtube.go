package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	youtube "github.com/kkdai/youtube/v2"
	"github.com/renato3x/mget/internal/utils"
	progressbar "github.com/schollz/progressbar/v3"
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
    label :=p.getFormatOptionLabel(f)
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

  for _, format := range formats {
    label := p.getFormatOptionLabel(format)
   
    err := p.download(&format, client, video, outputDir)
    
    if err != nil {
      fmt.Printf("Some error occurred downloading option %s\n", label)
    }
  }

  return nil
}

func (p *YouTubeProvider) createProgressBar(size int64, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(65 * time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func (p *YouTubeProvider) download(
  format *youtube.Format,
  client *youtube.Client,
  video *youtube.Video,
  outputDir string,
) error {
  stream, size, err := client.GetStream(video, format)
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

  label := p.getFormatOptionLabel(*format)
  desc := "Downloading "+label
  bar := p.createProgressBar(size, desc)
  _, err = io.Copy(io.MultiWriter(file, bar), stream)
  return err
}

func (p *YouTubeProvider) getFormatOptionLabel(format youtube.Format) string {
  if format.Height > 0 {
    return fmt.Sprintf("[%d] Audio + Video - %s (%s)", format.ItagNo, format.QualityLabel, format.MimeType)
  }

  return fmt.Sprintf("[%d] Audio only - (%s)", format.ItagNo, format.MimeType)
}
