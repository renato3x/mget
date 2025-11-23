# mget

A simple and efficient command-line tool for downloading media from various platforms. Currently supports YouTube with plans for additional platforms.

## Features

- 🎥 **Video Downloads**: Download videos in MP4 format
- 🎵 **Audio Extraction**: Download audio-only content in MP3 format
- 📊 **Progress Tracking**: Real-time download progress bars
- 🏠 **Organized Storage**: Automatically saves downloads to `~/mget-downloads`
- 🔗 **Multiple URL Formats**: Supports various YouTube URL formats (youtube.com, youtu.be, m.youtube.com)

## Installation

### Prerequisites

- Go 1.25.4 or later

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

Download a video from YouTube:

```bash
mget https://www.youtube.com/watch?v=VIDEO_ID
```

### Audio Only

Download only the audio track (MP3):

```bash
mget --audio https://www.youtube.com/watch?v=VIDEO_ID
# or use the short flag
mget -a https://www.youtube.com/watch?v=VIDEO_ID
```

### Supported URL Formats

The tool accepts various YouTube URL formats:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://m.youtube.com/watch?v=VIDEO_ID`

## How It Works

### Architecture

The project is organized into several packages:

- **`cmd/mget`**: Main entry point that handles CLI argument parsing and orchestrates the download process
- **`internal/cli`**: Command-line interface handling, including flag parsing
- **`internal/mget`**: Core download logic and platform-specific implementations

### Download Process

1. **URL Validation**: The tool validates and identifies the platform from the provided URL
2. **Platform Detection**: Currently supports YouTube
3. **Format Selection**: Automatically selects the best available format matching the requested media type (video or audio)
4. **Download**: Streams the media content with progress tracking
5. **File Storage**: Saves the file with a UUID-based filename to `~/mget-downloads`

### Platform Support

#### YouTube

- Downloads video in MP4 format
- Extracts audio in MP3 format
- Automatically selects formats with audio channels
- Uses the `kkdai/youtube/v2` library for YouTube API interaction

## Output Directory

All downloads are saved to:

```
~/mget-downloads/
```

Files are named using UUIDs to avoid conflicts. The directory is automatically created if it doesn't exist.

## Examples

### Download a video:

```bash
$ mget https://www.youtube.com/watch?v=dQw4w9WgXcQ
Downloading media from YouTube...
████████████████████████████████████████ 100%
Media "Rick Astley - Never Gonna Give You Up" downloaded to /Users/username/mget-downloads/550e8400-e29b-41d4-a716-446655440000.mp4
```

### Download audio only:

```bash
$ mget -a https://www.youtube.com/watch?v=dQw4w9WgXcQ
Downloading media from YouTube...
████████████████████████████████████████ 100%
Media "Rick Astley - Never Gonna Give You Up" downloaded to /Users/username/mget-downloads/550e8400-e29b-41d4-a716-446655440000.mp3
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
- `github.com/schollz/progressbar/v3`: Progress bar display
- `github.com/google/uuid`: UUID generation for filenames

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Roadmap

- [ ] Custom output directory option
- [ ] Quality/format selection
- [ ] Playlist support
- [ ] Additional platform support
