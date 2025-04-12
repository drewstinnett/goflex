package cmd

import (
	"log/slog"
	"strconv"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// deletePlaylistItemCmd represents the random command
var deletePlaylistItemCmd = &cobra.Command{
	Use:   "playlist-item PLAYLIST_NAME SHOW_NAME SEASON EPISODE",
	Short: "Delete an entry from a given playlist",
	Args:  cobra.MinimumNArgs(4),
	RunE: func(_ *cobra.Command, args []string) error {
		season, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		episode, err := strconv.Atoi(args[3])
		if err != nil {
			return err
		}
		p := newPlex()

		pl, err := p.Playlists.GetWithName(goflex.PlaylistTitle(args[0]))
		if err != nil {
			return err
		}
		slog.Info(
			"Removing item from playlist",
			"playlist",
			args[0],
			"show",
			args[0],
			"season",
			season,
			"episode",
			episode,
		)
		if err := p.Playlists.DeleteEpisode(pl.Title, goflex.ShowTitle(args[1]), goflex.SeasonNumber(season), goflex.EpisodeNumber(episode)); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deletePlaylistItemCmd)
	deletePlaylistItemCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
}
