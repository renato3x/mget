package main

import (
	"fmt"
	"os"

	"github.com/renato3x/mget/internal/cli"
)

func main() {
  args := cli.Args()

  if args.Url == "" {
    fmt.Println("URL is required")
    os.Exit(1)
    return
  }

  if err := cli.Download(args.Url, args.Output, args.AudioOnly); err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }

  os.Exit(0)
}
