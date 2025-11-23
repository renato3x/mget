package main

import (
	"fmt"
	"os"

	"github.com/renato3x/mget/internal/cli"
	"github.com/renato3x/mget/internal/mget"
)

func main() {
  args := cli.Args()

  if args.Url == "" {
    fmt.Println("URL is required")
    os.Exit(1)
    return
  }

  if err := mget.Download(args.Url, args.AudioOnly); err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }

  os.Exit(0)
}
