package cmd

import (
	"log/slog"

	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getPlaylistCmd represents the random command
var getPlaylistCmd = &cobra.Command{
	Use:   "playlist [TITLE]",
	Short: "Get a playlist from the API",
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		pl, err := p.Playlist(args[0])
		if err != nil {
			return err
		}
		episodes, err := pl.Episodes()
		if err != nil {
			return err
		}
		gout.MustPrint(episodes)
		slog.Info("completed", "episodes", len(episodes))
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPlaylistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	getPlaylistCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
	getPlaylistCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getPlaylistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
