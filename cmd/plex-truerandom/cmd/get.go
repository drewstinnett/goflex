package cmd

import (
	"context"
	"os"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
	"github.com/drewstinnett/gout/v2"
	plexrando "github.com/drewstinnett/plex-truerandom"
	"github.com/spf13/cobra"
)

// getCmd represents the random command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "getize a playlist",
	RunE: func(_ *cobra.Command, _ []string) error {
		p := plexrando.New(plexrando.WithAPI(plexgo.New(
			plexgo.WithSecurity(os.Getenv("PLEX_TOKEN")),
			plexgo.WithServerURL(os.Getenv("PLEX_URL")),
			plexgo.WithClientID("313FF6D7-5795-45E3-874F-B8FCBFD5E587"),
			plexgo.WithClientName("plex-trueget"),
			plexgo.WithClientVersion("0.0.1"),
		)))

		res, err := p.API.Playlists.GetPlaylistContents(context.Background(), p.PlaylistMap["American Dad!"].ID, operations.GetPlaylistContentsQueryParamTypeEpisode)
		if err != nil {
			return err
		}
		for _, episode := range res.Object.MediaContainer.Metadata {
			gout.MustPrint(episode)
			/*
				rk, err := strconv.ParseInt(*episode.RatingKey, 10, 64)
				if err != nil {
					return err
				}
				slog.Info("key", "key", rk)
				epi, err := p.API.Library.GetMetaDataByRatingKey(context.Background(), rk)
				if err != nil {
					panic(err)
					return err
				}
				for _, thing := range epi.Object.MediaContainer.Metadata {
					fmt.Fprintf(os.Stderr, "%+v\n", thing)
				}
			*/
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	getCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
	getCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
