package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getPlaylistEpisodesCmd represents the random command
var getPlaylistEpisodesCmd = &cobra.Command{
	Use:   "playlist-episodes TITLE",
	Short: "Get all of the episodes in a given playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		playlist, err := p.Playlists.GetWithName(args[0])
		if err != nil {
			return err
		}
		episodes, err := playlist.Episodes()
		if err != nil {
			return err
		}
		for _, item := range episodes {
			fmt.Println(item.String())
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPlaylistEpisodesCmd)
}
