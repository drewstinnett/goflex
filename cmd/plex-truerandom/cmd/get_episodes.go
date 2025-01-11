package cmd

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/spf13/cobra"
)

// getEspisodesCmd represents the random command
var getEpisodesCmd = &cobra.Command{
	Use:   "episodes SHOW",
	Short: "Get episodes",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		shows, err := p.MatchShows(args[0])
		if err != nil {
			return err
		}
		for _, show := range shows {
			slog.Info("show", "title", show.Title)
			seasons, err := show.Seasons()
			if err != nil {
				return err
			}
			seasonKeys := make([]int, len(seasons))
			i := 0
			for k := range seasons {
				seasonKeys[i] = k
				i++
			}
			sort.Ints(seasonKeys)
			for _, k := range seasonKeys {
				season := seasons[k]
				slog.Info("season", "show", show.Title, "index", season.Index, "key", season.ID)
				episodes, err := season.Episodes()
				if err != nil {
					return err
				}

				episodeKeys := make([]int, len(episodes))
				i := 0
				for k := range episodes {
					episodeKeys[i] = k
					i++
				}
				sort.Ints(episodeKeys)
				for _, k := range episodeKeys {
					episode := episodes[k]
					fmt.Println(episode.String())
				}
			}
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getEpisodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getEspisodesCmd.PersistentFlags().String("library", "TV Seasons", "Library of the TV Show we are randomizing")
	// getEspisodesCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getEspisodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
