package cmd

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getEspisodesCmd represents the random command
var getEpisodesCmd = &cobra.Command{
	Use:   "episodes SHOW",
	Short: "Get episodes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		shows, err := p.MatchShows(args[0])
		if err != nil {
			return err
		}
		short := mustGetCmd[bool](*cmd, "short")
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
				if short {
					for _, k := range episodeKeys {
						episode := episodes[k]
						fmt.Println(episode.String())
					}
				} else {
					gout.MustPrint(episodes)
				}
			}
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getEpisodesCmd)
	getEpisodesCmd.PersistentFlags().BoolP("short", "s", false, "Show short version of the episode (Name S00E00)")
}
