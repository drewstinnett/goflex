package cmd

import (
	"github.com/spf13/cobra"
)

// markWatchedCmd represents the random command
var markWatchedCmd = &cobra.Command{
	Use:   "mark-watched SHOW SEASON EPISODE",
	Short: "markWatched something from the plex server",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()
		_ = p
		return nil
	},
}

func init() {
	rootCmd.AddCommand(markWatchedCmd)
}
