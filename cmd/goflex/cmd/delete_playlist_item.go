package cmd

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
)

// deletePlaylistItemCmd represents the random command
var deletePlaylistItemCmd = &cobra.Command{
	Use:   "playlist-item PLAYLIST_NAME SHOW_NAME SEASON EPISODE",
	Short: "Delete an entry from a given playlist",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(_ *cobra.Command, args []string) error {
		slog.Info("DING", "args", args)
		season, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		episode, err := strconv.Atoi(args[3])
		if err != nil {
			return err
		}
		p := newPlex()

		pl, err := p.Playlist(args[0])
		if err != nil {
			return err
		}
		_ = season
		_ = episode
		_ = pl
		/*
			slog.Info("Removing item from playlist", "playlist", args[0], "season", season, "episode", episode)
			if err := pl.DeleteEpisode(args[1], season, episode); err != nil {
				return err
			}
		*/
		return errors.New("not yet implemented, see above")
	},
}

func init() {
	deleteCmd.AddCommand(deletePlaylistItemCmd)
	deletePlaylistItemCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
}
