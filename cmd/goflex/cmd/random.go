package cmd

import (
	"log/slog"
	"math/rand"
	"time"

	plexrando "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random PLAYLIST_NAME",
	Short: "Randomize a playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		// Inspect the playlist, create it if it doesn't exist
		playlist, created, err := p.Playlists.GetOrCreate(args[0], plexrando.VideoPlaylist, false)
		if err != nil {
			return err
		}

		var playlistEpisodes plexrando.EpisodeList
		var refillPlaylist bool
		if created {
			refillPlaylist = true
		} else {
			var err error
			playlistEpisodes, err = playlist.Episodes()
			if err != nil {
				return err
			}
			slog.Debug("found existing playlist", "count", len(playlistEpisodes))
		}
		showTitle := mustGetCmd[string](*cmd, "title")

		// Get viewed
		viewed, err := p.Sessions.HistoryEpisodes(toPTR(time.Now().Add(-time.Hour*24*time.Duration(mustGetCmd[int](*cmd, "lookback-days")))), showTitle)
		if err != nil {
			return err
		}
		slog.Debug("found viewed episodes", "count", len(viewed))

		// Figure out remaining
		remaining, removed := playlistEpisodes.Subtract(viewed)
		if len(remaining) <= mustGetCmd[int](*cmd, "refill-at") {
			refillPlaylist = true
		}
		slog.Debug("removed viewed episodes", "removed", len(removed), "remaining", len(remaining))

		// Remove things we have seen
		if len(removed) > 0 {
			slog.Info("New length of episodes after removing viewed", "remaining", len(remaining), "removed", len(removed), "original", len(playlistEpisodes))
			for _, item := range removed {
				slog.Info("removing episode", "playlist", args[0], "episode", item.String())
				if err := p.Playlists.DeleteEpisode(playlist.Title, item.Show, item.Season, item.Episode); err != nil {
					return err
				}
			}
		}

		if refillPlaylist {
			slog.Info("refilling playlist", "playlist", args[0])
			if err := p.Playlists.Clear(playlist.ID); err != nil {
				return err
			}

			shows, err := p.MatchShows(mustGetCmd[string](*cmd, "title"))
			if err != nil {
				return err
			}

			allEpisodes, err := shows.Episodes()
			if err != nil {
				return err
			}

			unviewedEpisodes, _ := allEpisodes.Subtract(viewed)
			rand.Shuffle(len(unviewedEpisodes), func(i, j int) {
				unviewedEpisodes[i], unviewedEpisodes[j] = unviewedEpisodes[j], unviewedEpisodes[i]
			})
			slog.Debug("refilling playlist", "title", playlist.Title, "episodes", len(unviewedEpisodes))
			return p.Playlists.InsertEpisodes(playlist.ID, unviewedEpisodes)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	/*
		randomCmd.PersistentFlags().String("library", "", "Library of the TV Show we are randomizing")
		if err := randomCmd.MarkPersistentFlagRequired("library"); err != nil {
			panic(err)
		}
	*/
	randomCmd.PersistentFlags().String("title", "", "Name of the show to include in this playlist")
	if err := randomCmd.MarkPersistentFlagRequired("title"); err != nil {
		panic(err)
	}

	randomCmd.PersistentFlags().Int("lookback-days", 14, "number of days to look back at viewed history")
	randomCmd.PersistentFlags().Int("refill-at", 10, "refill the playlist when it reaches this remaining number of episodes")
	randomCmd.PersistentFlags().Int("earliest-season", 0, "earliest season to include in the playlist")
	randomCmd.PersistentFlags().Int("latest-season", 0, "latest season to include in the playlist")
}
