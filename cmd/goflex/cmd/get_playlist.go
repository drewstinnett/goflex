package cmd

import (
	goflex "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
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
		ret := make([]*goflex.Playlist, len(args))
		for idx, item := range args {
			got, err := p.Playlists.GetWithName(goflex.PlaylistTitle(item))
			if err != nil {
				return err
			}
			ret[idx] = got
		}
		gout.MustPrint(ret)

		return nil
	},
}

func init() {
	getCmd.AddCommand(getPlaylistCmd)
}
