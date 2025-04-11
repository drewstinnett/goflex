package cmd

import (
	"errors"
	"log/slog"
	"time"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// getEspisodesCmd represents the random command
var getEpisodesCmd = &cobra.Command{
	Use:   "episodes SHOW",
	Short: "Get episodes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		shows, err := p.Shows.Match(goflex.ShowTitle(args[0]))
		if err != nil {
			return err
		}
		if len(shows) == 0 {
			return errors.New("show not found: " + args[0])
		}
		short := mustGetCmd[bool](*cmd, "short")
		wd := mustGetCmd[time.Duration](*cmd, "watch-duration")
		var watch bool
		if wd > 0 {
			watch = true
		}
		for {
			for _, show := range shows {
				slog.Info("show", "title", show.Title)
				seasons, err := p.Shows.SeasonsSorted(*show)
				if err != nil {
					return err
				}
				for _, season := range seasons {
					slog.Debug("season", "show", show.Title, "index", season.Index, "key", season.ID)
					// episodesM, err := season.Episodes()
					episodesM, err := p.Shows.SeasonEpisodes(season)
					if err != nil {
						return err
					}
					printEpisodes(episodesM.List(), short)

				}
			}
			if !watch {
				return nil
			}
			time.Sleep(wd)
		}
	},
}

func init() {
	getEpisodesCmd.PersistentFlags().DurationP("watch-duration", "w", 0, "get seasons again every duration")
	getEpisodesCmd.PersistentFlags().BoolP("short", "s", false, "Show short version of the episode (Name S00E00)")
	getCmd.AddCommand(getEpisodesCmd)
}
