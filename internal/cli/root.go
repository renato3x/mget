package cli

import (
	"fmt"
	"os"

	"github.com/renato3x/mget/internal/downloader"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
  Use: "mget [url]",
  Short: "A simple and efficient command-line tool for downloading media from various platforms",
  Args: cobra.MinimumNArgs(1),
  RunE: func (cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      cmd.Help()
      return nil
    }
    url := args[0]
    return downloader.Download(url)
  },
}

func Start() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
