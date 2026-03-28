# mget

A simple and efficient command-line tool for downloading media from various platforms. Currently supports YouTube with plans for additional platforms.

## Features

- �️ **Interactive Format Selection**: Choose one or more formats to download via an interactive prompt
- 📊 **Parallel Downloads**: Selected formats are downloaded simultaneously, each with its own progress bar
- 🏠 **Organized Storage**: Automatically saves downloads to `~/mget-downloads`
- 🔗 **Multiple URL Formats**: Supports various YouTube URL formats (youtube.com, youtu.be, m.youtube.com, Shorts, Embed)

## Installation

### Prerequisites

- Go 1.26 or later

### Build from Source

```bash
git clone https://github.com/renato3x/mget.git
cd mget
go build -o mget ./cmd/mget
```

### Install with Go

Install directly from GitHub:

```bash
go install github.com/renato3x/mget/cmd/mget@latest
```

Or install from a local clone:

```bash
go install ./cmd/mget
```

Make sure your `$GOPATH/bin` or `$GOBIN` is in your `$PATH`.

## Usage

### Basic Usage

Provide a URL and follow the interactive prompt to select the formats you want to download:

```bash
mget https://www.youtube.com/watch?v=VIDEO_ID
```

### Check Version

```bash
mget version
```

### Supported URL Formats

The tool accepts various YouTube URL formats:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://m.youtube.com/watch?v=VIDEO_ID`
- `https://www.youtube.com/shorts/VIDEO_ID`
- `https://www.youtube.com/embed/VIDEO_ID`
- `https://www.youtube.com/v/VIDEO_ID`

## How It Works

### Architecture

The project is organized into several packages:

- **`cmd/mget`**: Main entry point
- **`internal/cli`**: Command-line interface built with `cobra`, including the root command and subcommands
- **`internal/downloader`**: Core download logic, provider interface, and platform-specific implementations

### Download Process

1. **URL Validation**: The tool validates the provided URL
2. **Platform Detection**: Matches the URL against registered providers (currently YouTube)
3. **Format Listing**: Fetches all available formats with audio channels for the requested media
4. **Interactive Selection**: Presents a multi-select prompt so you can pick one or more formats to download
5. **Parallel Download**: Streams all selected formats concurrently, each with its own progress bar
6. **File Storage**: Saves each file to `~/mget-downloads` with a slugified name derived from the format tag, MIME type, and media title

### Platform Support

#### YouTube

- Supports audio-only and audio+video formats
- Lists all available formats with audio channels for selection
- Uses the `kkdai/youtube/v2` library for YouTube API interaction

## Output Directory

All downloads are saved to:

```
~/mget-downloads/
```

Files are named using a slug built from the format's itag number, MIME type, and the media title (e.g., `137-video-mp4-rick-astley-never-gonna-give-you-up.mp4`). The directory is automatically created if it doesn't exist.

## Examples

### Download a video:

```bash
$ mget https://www.youtube.com/watch?v=dQw4w9WgXcQ
Rick Astley - Never Gonna Give You Up - 3m33s

? Select one or more options for download  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
> [ ]  [137] Audio + Video - 1080p (video/mp4; codecs="avc1.640028")
  [ ]  [248] Audio + Video - 1080p (video/webm; codecs="vp9")
  [ ]  [140] Audio only - (audio/mp4; codecs="mp4a.40.2")
  ...
```

### Check the current version:

```bash
$ mget version
0.4.0
```

## Error Handling

The tool provides clear error messages for common issues:

- Invalid or missing URL
- Unsupported platform
- Network errors
- Format selection failures
- File system errors

## Dependencies

- `github.com/kkdai/youtube/v2`: YouTube video downloading
- `github.com/AlecAivazis/survey/v2`: Interactive multi-select prompt
- `github.com/vbauerster/mpb/v8`: Multi-bar concurrent progress display
- `github.com/spf13/cobra`: CLI framework

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.
