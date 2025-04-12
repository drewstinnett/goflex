package cmd

import (
	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// getPlaylistEpisodesCmd represents the random command
var getPlaylistEpisodesCmd = &cobra.Command{
	Use:   "playlist-episodes TITLE",
	Short: "Get all of the episodes in a given playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		playlist, err := p.Playlists.GetWithName(goflex.PlaylistTitle(args[0]))
		if err != nil {
			return err
		}
		// episodes, err := playlist.Episodes()
		episodes, err := p.Playlists.Episodes(*playlist)
		if err != nil {
			return err
		}

		printEpisodes(episodes, mustGetCmd[bool](*cmd, "short"))

		return nil
	},
}

func init() {
	getCmd.AddCommand(getPlaylistEpisodesCmd)
	getPlaylistEpisodesCmd.PersistentFlags().
		BoolP("short", "s", false, "Show short version of the episode (Name S00E00)")
}
