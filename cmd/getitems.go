// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/choppedpork/ohctl/openhab"

	"github.com/spf13/cobra"
)

var quiet bool

// getitemsCmd represents the get items command
var getitemsCmd = &cobra.Command{
	Use:   "items",
	Short: "List all items",
	Long:  `List all items in openhab`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("get items called")

		oh := openhab.NewClient()
		items := oh.GetItems()

		if quiet {
			for _, item := range items {
				fmt.Println(item.Name)
			}
		} else {
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 12, 8, 1, '\t', 0)

			fmt.Fprintf(w, "name\tstate\tgroups\ttags\n")
			fmt.Fprintf(w, "----\t-----\t------\t----\n")
			for _, item := range items {
				fmt.Fprintf(w, "%s\t%.12s\t(%v)\t[%v]\n", item.Name, item.State,
					strings.Join(item.GroupNames, ", "), strings.Join(item.Tags, ", "))
			}

			w.Flush()
		}
	},
}

func init() {
	getitemsCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode - prints item names only")
	getCmd.AddCommand(getitemsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getitemsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getitemsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
