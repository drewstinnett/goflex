package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// randomCmdDeprecated represents the random command
var randomCmdDeprecated = &cobra.Command{
	Use:   "randomdeprecated PLAYLIST_NAME",
	Short: "Randomize a playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		// Inspect the playlist, create it if it doesn't exist
		playlist, created, err := p.Playlists.GetOrCreate(goflex.PlaylistTitle(args[0]), goflex.VideoPlaylist, false)
		if err != nil {
			return err
		}

		var playlistEpisodes goflex.EpisodeList
		var refillPlaylist bool
		var refillReason string
		if created {
			refillPlaylist = true
			refillReason = "newly created playlist"
		} else {
			var err error
			// playlistEpisodes, err = playlist.Episodes()
			playlistEpisodes, err = p.Playlists.Episodes(*playlist)
			if err != nil {
				return err
			}
			slog.Debug("found existing playlist", "count", len(playlistEpisodes))
		}
		showTitle := goflex.ShowTitle(mustGetCmd[string](*cmd, "title"))

		exists, err := p.Shows.Exists(showTitle)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("show does not exist: " + string(showTitle))
		}

		// Get viewed
		viewed, err := p.Sessions.HistoryEpisodes(
			time.Now().Add(-time.Hour*24*time.Duration(mustGetCmd[int](*cmd, "lookback-days"))),
			showTitle,
		)
		if err != nil {
			return err
		}
		slog.Debug("found viewed episodes", "count", len(viewed))

		// Figure out remaining
		remaining, removed := playlistEpisodes.Subtract(viewed)
		refillAt := mustGetCmd[int](*cmd, "refill-at")
		if (!refillPlaylist) && len(remaining) <= refillAt {
			refillPlaylist = true
			refillReason = fmt.Sprintf("playlist dipped below refill-at level: %v", refillAt)
		}
		slog.Debug("removed viewed episodes", "removed", len(removed), "remaining", len(remaining))

		// Remove things we have seen
		if len(removed) > 0 {
			slog.Debug(
				"New length of episodes after removing viewed",
				"remaining",
				len(remaining),
				"removed",
				len(removed),
				"original",
				len(playlistEpisodes),
			)
			for _, item := range removed {
				slog.Info("removing episode", "playlist", args[0], "episode", item.String())
				if err := p.Playlists.DeleteEpisode(playlist.Title, item.Show, item.Season, item.Episode); err != nil {
					return err
				}
			}
		}

		if refillPlaylist {
			slog.Debug("attempting to refill playlist", "playlist", args[0], "reason", refillReason)
			if err := p.Playlists.Clear(*playlist); err != nil {
				return err
			}

			// title := mustGetCmd[string](*cmd, "title")
			shows, err := p.Shows.Match(showTitle)
			if err != nil {
				return err
			}

			// allEpisodes, err := shows.EpisodesWithFilter(goflex.EpisodeFilter{
			allEpisodes, err := p.Shows.EpisodesWithFilter(shows, goflex.EpisodeFilter{
				LatestSeason:   goflex.SeasonNumber(mustGetCmd[int](*cmd, "latest-season")),
				EarliestSeason: goflex.SeasonNumber(mustGetCmd[int](*cmd, "earliest-season")),
			})
			if err != nil {
				return err
			}

			unviewedEpisodes, _ := allEpisodes.Subtract(viewed)
			rand.Shuffle(len(unviewedEpisodes), func(i, j int) {
				unviewedEpisodes[i], unviewedEpisodes[j] = unviewedEpisodes[j], unviewedEpisodes[i]
			})

			if len(unviewedEpisodes) < refillAt {
				return fmt.Errorf(
					"not enough unwatched episodes to refill. unwatched: %v, refill-at: %v",
					len(unviewedEpisodes),
					refillAt,
				)
			} else {
				slog.Info("refilling playlist", "title", playlist.Title, "episodes", len(unviewedEpisodes), "reason", refillReason)
				return p.Playlists.InsertEpisodes(playlist.ID, unviewedEpisodes)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randomCmdDeprecated)
	randomCmdDeprecated.PersistentFlags().Int("lookback-days", 14, "number of days to look back at viewed history")
	randomCmdDeprecated.PersistentFlags().
		Int("refill-at", 10, "refill the playlist when it reaches this remaining number of episodes")

	if err := bindShowFilter(randomCmdDeprecated); err != nil {
		panic(err)
	}
}

func bindShowFilter(cmd *cobra.Command) error {
	cmd.PersistentFlags().String("title", "", "Name of the show to include in this playlist")
	if err := cmd.MarkPersistentFlagRequired("title"); err != nil {
		return err
	}

	cmd.PersistentFlags().Int("earliest-season", 0, "earliest season to include in the playlist")
	cmd.PersistentFlags().Int("latest-season", 0, "latest season to include in the playlist")
	return nil
}
