package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " " + version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
