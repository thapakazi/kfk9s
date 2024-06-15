// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "kfk9s",
    Short: "kfk9s is a CLI tool for managing Kafka clusters",
}

func Execute() error {
    return rootCmd.Execute()
}
