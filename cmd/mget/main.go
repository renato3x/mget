package main

import (
	"fmt"
	"os"

	"github.com/renato3x/mget/internal/cli"
	"github.com/renato3x/mget/internal/mget"
)

const version = "0.3.2"

func main() {
  args := cli.Args()

  if args.Version {
    fmt.Printf("mget version %s\n", version)
    os.Exit(0)
    return
  }

  if args.Url == "" {
    showUsage()
    os.Exit(0)
    return
  }

  if err := mget.Download(args.Url, args.AudioOnly); err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }

  os.Exit(0)
}

func showUsage() {
  fmt.Println("mget - A simple and efficient command-line tool for downloading media from various platforms")
  fmt.Println()
  fmt.Println("Usage:")
  fmt.Println("  mget [flags] <URL>")
  fmt.Println()
  fmt.Println("Flags:")
  fmt.Println("  -a, --audio    Download only the audio track (MP3 format)")
  fmt.Println("  -v, --version  Show version information")
  fmt.Println()
  fmt.Println("Examples:")
  fmt.Println("  mget https://www.youtube.com/watch?v=VIDEO_ID")
  fmt.Println("  mget --audio https://www.youtube.com/watch?v=VIDEO_ID")
  fmt.Println("  mget -a https://youtu.be/VIDEO_ID")
  fmt.Println()
  fmt.Println("Supported Platforms:")
  fmt.Println("  - YouTube (youtube.com, youtu.be, m.youtube.com)")
  fmt.Println()
  fmt.Println("Output:")
  fmt.Println("  Downloads are saved to ~/mget-downloads/")
  fmt.Println("  Files are automatically named with UUIDs to avoid conflicts")
}
