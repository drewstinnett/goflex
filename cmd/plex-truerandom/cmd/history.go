package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// historyCmd represents the random command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Session history",
	RunE: func(cmd *cobra.Command, _ []string) error {
		p := newPlex()

		viewed, err := p.Viewed(mustGetCmd[string](*cmd, "library"), time.Now().Add(-time.Hour*24*14))
		if err != nil {
			return err
		}

		for _, episode := range viewed {
			fmt.Println(episode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
}
