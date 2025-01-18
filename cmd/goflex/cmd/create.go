package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// getCmd represents the random command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create something from the plex server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, _ []string) error {
		return errors.New(cmd.UsageString())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
