package cmd

import (
	"log/slog"
	"time"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getShowsCmd represents the random command
var getShowsCmd = &cobra.Command{
	Use:     "shows",
	Short:   "Get shows",
	Aliases: []string{"show"},
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		p := newPlex()

		libs, err := p.Library.List()
		if err != nil {
			return err
		}
		wd := mustGetCmd[time.Duration](*cmd, "watch-duration")
		var watch bool
		if wd > 0 {
			watch = true
		}
		for {
			for _, lib := range libs {
				if lib.Type != goflex.ShowType {
					continue
				}
				slog.Info("shows in library", "library", lib.Title)
				shows, err := p.Library.Shows(*lib)
				if err != nil {
					return err
				}
				gout.MustPrint(shows)
			}
			if !watch {
				return nil
			}
			time.Sleep(wd)
		}
	},
}

func init() {
	getShowsCmd.PersistentFlags().DurationP("watch-duration", "w", 0, "get seasons again every duration")
	getCmd.AddCommand(getShowsCmd)
}
