package cmd

import (
	"log/slog"

	plexrando "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getShowsCmd represents the random command
var getShowsCmd = &cobra.Command{
	Use:   "shows",
	Short: "Get shows",
	Args:  cobra.ExactArgs(0),
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		libs, err := p.Libraries()
		if err != nil {
			return err
		}
		for _, lib := range libs {
			if lib.Type != plexrando.ShowType {
				continue
			}
			slog.Info("shows in library", "library", lib.Title)
			shows, err := lib.Shows()
			if err != nil {
				return err
			}
			gout.MustPrint(shows)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getShowsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getShowsCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
	// getShowsCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getShowsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
