package cmd

import (
	"errors"
	"log/slog"

	"github.com/spf13/cobra"
)

// getEspisodesCmd represents the random command
var getEpisodesCmd = &cobra.Command{
	Use:   "episodes SHOW",
	Short: "Get episodes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		shows, err := p.Shows.Match(args[0])
		if err != nil {
			return err
		}
		if len(shows) == 0 {
			return errors.New("show not found: " + args[0])
		}
		short := mustGetCmd[bool](*cmd, "short")
		for _, show := range shows {
			slog.Info("show", "title", show.Title)
			seasons, err := show.SeasonsSorted()
			if err != nil {
				return err
			}
			for _, season := range seasons {
				slog.Debug("season", "show", show.Title, "index", season.Index, "key", season.ID)
				episodesM, err := season.Episodes()
				if err != nil {
					return err
				}
				printEpisodes(episodesM.List(), short)

			}
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getEpisodesCmd)
	getEpisodesCmd.PersistentFlags().BoolP("short", "s", false, "Show short version of the episode (Name S00E00)")
}
