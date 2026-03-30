package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	survey "github.com/AlecAivazis/survey/v2"
	youtube "github.com/kkdai/youtube/v2"
	"github.com/renato3x/mget/internal/utils"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type YouTubeProvider struct{}

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

	fmt.Println()
	return p.downloadOptions(options, &client, video)
}

func (p *YouTubeProvider) selectOptions(formats youtube.FormatList) ([]youtube.Format, error) {
	var options []string
	formatMap := make(map[string]*youtube.Format)

	for _, f := range formats {
		label := p.getFormatOptionLabel(f)
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

	pContainer := mpb.New(mpb.WithWidth(60))
	var wg sync.WaitGroup

	for _, format := range formats {
		wg.Add(1)
		label := p.getFormatOptionLabel(format)
		go func(f youtube.Format) {
			defer wg.Done()
			err := p.download(&f, client, video, outputDir, pContainer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Some error occurred downloading option %s\n", label)
			}
		}(format)
	}

	wg.Wait()
	pContainer.Wait()
	fmt.Println("\nAll downloads were done.")
	return nil
}

func (p *YouTubeProvider) download(
	format *youtube.Format,
	client *youtube.Client,
	video *youtube.Video,
	outputDir string,
	container *mpb.Progress,
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
	bar := container.AddBar(
		size,
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("Downloading %s", label)),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.OnComplete(decor.Name(" ✓"), " done"),
		),
	)

	proxyReader := bar.ProxyReader(stream)
	defer proxyReader.Close()

	_, err = io.Copy(file, proxyReader)
	return err
}

func (p *YouTubeProvider) getFormatOptionLabel(format youtube.Format) string {
  ext := strings.ToUpper(utils.GetFileExtensionByMimetype(format.MimeType)[1:])
	if format.Height > 0 {
		return fmt.Sprintf("[%d] %s (%s)", format.ItagNo, format.QualityLabel, ext)
	}

	return fmt.Sprintf("[%d] Audio (%s)", format.ItagNo, ext)
}
