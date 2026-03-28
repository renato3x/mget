package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows current mget version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.3.2")
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
