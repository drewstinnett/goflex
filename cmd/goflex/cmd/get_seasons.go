package cmd

import (
	"log/slog"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getSeasonsCmd represents the random command
var getSeasonsCmd = &cobra.Command{
	Use:   "seasons SHOW",
	Short: "Get shows",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		shows, err := p.Shows.Match(goflex.ShowTitle(args[0]))
		if err != nil {
			return err
		}
		for _, show := range shows {
			slog.Info("show", "title", show.Title)
			seasons, err := p.Shows.Seasons(*show)
			if err != nil {
				return err
			}
			gout.MustPrint(seasons)
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getSeasonsCmd)
}
