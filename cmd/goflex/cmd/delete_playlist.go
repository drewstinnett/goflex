package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

// deletePlaylistItemCmd represents the random command
var deletePlaylistCmd = &cobra.Command{
	Use:   "playlist PLAYLIST_NAME",
	Short: "Delete a playlist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		pl, err := p.Playlists.GetWithName(args[0])
		if err != nil {
			return err
		}
		slog.Info("deleting", "playlist", pl.Title, "id", pl.ID)
		return p.Playlists.Delete(pl.ID)
	},
}

func init() {
	deleteCmd.AddCommand(deletePlaylistCmd)
}
