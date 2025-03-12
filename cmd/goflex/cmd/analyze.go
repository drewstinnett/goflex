package cmd

import (
	"fmt"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// getCmd represents the random command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze something from the plex server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		p := newPlex()
		shows, err := p.Shows.StrictMatch(goflex.ShowTitle(mustGetCmd[string](*cmd, "title")))
		if err != nil {
			return err
		}

		// allEpisodes, err := shows.EpisodesWithFilter(goflex.EpisodeFilter{
		allEpisodes, err := p.Shows.EpisodesWithFilter(shows, goflex.EpisodeFilter{
			LatestSeason:   mustGetCmd[int](*cmd, "latest-season"),
			EarliestSeason: mustGetCmd[int](*cmd, "earliest-season"),
		})
		if err != nil {
			return err
		}

		duration, err := allEpisodes.WatchSpan()
		if err != nil {
			return err
		}

		fmt.Printf("Took %v days to watch\n", int(duration.Hours()/24))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	if err := bindShowFilter(analyzeCmd); err != nil {
		panic(err)
	}
}
