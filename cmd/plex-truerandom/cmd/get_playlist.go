package cmd

import (
	"github.com/drewstinnett/gout/v2"
	plexrando "github.com/drewstinnett/plex-truerandom"
	"github.com/spf13/cobra"
)

// getPlaylistCmd represents the random command
var getPlaylistCmd = &cobra.Command{
	Use:   "playlist [TITLE]",
	Short: "Get a playlist from the API",
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		if len(args) == 0 {
			ret, err := p.Playlists.List()
			if err != nil {
				return err
			}
			return gout.Print(ret)
		}
		ret := make([]*plexrando.Playlist, len(args))
		for idx, item := range args {
			got, err := p.Playlists.GetWithName(item)
			if err != nil {
				return err
			}
			ret[idx] = got
		}
		gout.MustPrint(ret)

		/*
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
		*/
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPlaylistCmd)
}
