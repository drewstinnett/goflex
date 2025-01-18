package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// markWatchedCmd represents the random command
var markUnWatchedCmd = &cobra.Command{
	Use:     "mark-unwatched SHOW SEASON EPISODE",
	Short:   "markWatched something from the plex server",
	Aliases: []string{"unwatch", "unwatched"},
	Args:    cobra.ExactArgs(3),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()
		show, season, episode, err := episodeArgs(args)
		if err != nil {
			return err
		}
		if err := p.Media.MarkEpisodeUnWatched(show, season, episode); err != nil {
			return err
		}
		fmt.Println("Marked episode as un-watched!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(markUnWatchedCmd)
}
