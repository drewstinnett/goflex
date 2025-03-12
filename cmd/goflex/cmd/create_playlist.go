package cmd

import (
	goflex "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// createPlaylistCmd represents the random command
var createPlaylistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Create a new playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()
		if err := p.Playlists.Create(goflex.PlaylistTitle(args[0]), "video", false); err != nil {
			return err
		}
		playlist, err := p.Playlists.GetWithName(goflex.PlaylistTitle(args[0]))
		if err != nil {
			return err
		}
		return gout.Print(playlist)
	},
}

func init() {
	createCmd.AddCommand(createPlaylistCmd)
}
