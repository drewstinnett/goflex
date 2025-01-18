package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// markWatchedCmd represents the random command
var markWatchedCmd = &cobra.Command{
	Use:     "mark-watched SHOW SEASON EPISODE",
	Short:   "Mark something as watched on Plex",
	Aliases: []string{"watch", "watched"},
	Args:    cobra.ExactArgs(3),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		show, season, episode, err := episodeArgs(args)
		if err != nil {
			return err
		}
		if err := p.Media.MarkEpisodeWatched(show, season, episode); err != nil {
			return err
		}
		fmt.Println("Marked episode as watched!")
		return nil
	},
}

func episodeArgs(args []string) (string, int, int, error) {
	show := args[0]
	seasonRaw := args[1]
	episodeRaw := args[2]
	season, err := strconv.Atoi(seasonRaw)
	if err != nil {
		return "", 0, 0, err
	}
	episode, err := strconv.Atoi(episodeRaw)
	if err != nil {
		return "", 0, 0, err
	}
	return show, season, episode, nil
}

func init() {
	rootCmd.AddCommand(markWatchedCmd)
}
