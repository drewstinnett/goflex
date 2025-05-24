package cmd

import (
	"log/slog"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// deletePlaylistItemCmd represents the random command
var deletePlaylistCmd = &cobra.Command{
	Use:   "playlist PLAYLIST_NAME",
	Short: "Delete a playlist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		pl, err := p.Playlists.GetWithName(goflex.PlaylistTitle(args[0]))
		if err != nil {
			return err
		}
		slog.Info("deleting", "playlist", pl.Title, "id", pl.ID)
		return p.Playlists.Delete(*pl)
	},
}

func init() {
	deleteCmd.AddCommand(deletePlaylistCmd)
}
