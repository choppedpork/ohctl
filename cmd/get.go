package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:       "get",
	Short:     "Gets stuff",
	ValidArgs: []string{"item", "items"},
	// no Run here
}

func init() {
	RootCmd.AddCommand(getCmd)
}
