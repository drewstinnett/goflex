package cmd

import (
	"log/slog"

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

		shows, err := p.Shows.Match(args[0])
		if err != nil {
			return err
		}
		for _, show := range shows {
			slog.Info("show", "title", show.Title)
			seasons, err := show.Seasons()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getSeasonsCmd.PersistentFlags().String("library", "TV Seasons", "Library of the TV Show we are randomizing")
	// getSeasonsCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getSeasonsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
