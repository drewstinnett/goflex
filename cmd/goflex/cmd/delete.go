package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// deleteCmd represents the random command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete something from the plex server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return errors.New(cmd.UsageString())
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
