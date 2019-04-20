package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/choppedpork/ohctl/openhab"
	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd <item> <command>",
	Short: "Sends a command to an item.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		oh := openhab.NewClient(Config.Host, Config.Port)
		err := oh.Cmd(args[0], strings.ToUpper(args[1]))

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(cmdCmd)
}
