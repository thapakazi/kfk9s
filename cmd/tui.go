// cmd/tui.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thapakazi/kfk9s/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the TUI",
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartUI()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
