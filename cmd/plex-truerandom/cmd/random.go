package cmd

import (
	"log/slog"
	"math/rand"
	"time"

	plexrando "github.com/drewstinnett/plex-truerandom"
	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random NEW_PLAYLIST_NAME",
	Short: "Randomize a playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		playlist, created, err := p.GetOrCreatePlaylist(args[0])
		if err != nil {
			return err
		}

		// Shuffle if new
		var episodes plexrando.EpisodeList
		var allEpisodes plexrando.EpisodeList
		if created {
			// Find the stuff currently in the playlist
			var err error
			allEpisodes, err = p.Episodes(
				mustGetCmd[string](*cmd, "library"),
				mustGetCmd[string](*cmd, "title"),
			)
			if err != nil {
				return err
			}
			allEpisodes = allEpisodes.Seasons(mustGetCmd[int](*cmd, "earliest-season"), mustGetCmd[int](*cmd, "latest-season"))
			slog.Debug("found existing episodes", "count", len(allEpisodes))
			slog.Info("shuffling because we just created this playlist")
			rand.Shuffle(len(allEpisodes), func(i, j int) { allEpisodes[i], allEpisodes[j] = allEpisodes[j], allEpisodes[i] })
			episodes = allEpisodes
		} else {
			playlistEpisodes, err := p.PlaylistEpisodes(args[0])
			if err != nil {
				return err
			}
			slog.Debug("found existing playlist", "count", len(playlistEpisodes))
			episodes = playlistEpisodes

		}

		// Get viewed
		viewed, err := p.Viewed(mustGetCmd[string](*cmd, "library"), time.Now().Add(-time.Hour*24*time.Duration(mustGetCmd[int](*cmd, "lookback-days"))))
		if err != nil {
			return err
		}
		slog.Debug("found viewed episodes", "count", len(viewed))

		var rewritePlaylist bool
		// Figure out remaining
		remaining, removed := episodes.Subtract(viewed)
		if len(removed) > 0 {
			slog.Info("New length of episodes after removing viewed", "remaining", len(remaining), "removed", len(removed), "original", len(episodes))
			rewritePlaylist = true
			for _, item := range removed {
				slog.Info("removing episode", "episode", item.String())
			}
		}

		if len(remaining) <= mustGetCmd[int](*cmd, "refill-at") {
			slog.Info("playlist reached low level, re-shuffling", "remaining", len(remaining))
			rewritePlaylist = true
			allEpisodes, err = p.Episodes(
				mustGetCmd[string](*cmd, "library"),
				mustGetCmd[string](*cmd, "title"),
			)
			if err != nil {
				return err
			}
			allEpisodes = allEpisodes.Seasons(mustGetCmd[int](*cmd, "earliest-season"), mustGetCmd[int](*cmd, "latest-season"))
			// Refill remaining
			remaining = allEpisodes
			rand.Shuffle(len(remaining), func(i, j int) { remaining[i], remaining[j] = remaining[j], remaining[i] })
		}

		if rewritePlaylist {
			if err := playlist.Clear(); err != nil {
				return err
			}

			return playlist.AddEpisodes(remaining)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	randomCmd.PersistentFlags().String("library", "", "Library of the TV Show we are randomizing")
	if err := randomCmd.MarkPersistentFlagRequired("library"); err != nil {
		panic(err)
	}
	randomCmd.PersistentFlags().String("title", "", "Name of the show to include in this playlist")
	if err := randomCmd.MarkPersistentFlagRequired("title"); err != nil {
		panic(err)
	}

	randomCmd.PersistentFlags().Int("lookback-days", 14, "number of days to look back at viewed history")
	randomCmd.PersistentFlags().Int("refill-at", 10, "refill the playlist when it reaches this remaining number of episodes")
	randomCmd.PersistentFlags().Int("earliest-season", 0, "earliest season to include in the playlist")
	randomCmd.PersistentFlags().Int("latest-season", 0, "latest season to include in the playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
