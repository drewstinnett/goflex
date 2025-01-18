package cmd

import (
	"fmt"
	"time"

	plexrando "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

// getSessionsCmd represents the random command
var getSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Get session information",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		var ret plexrando.EpisodeList
		if mustGetCmd[bool](*cmd, "history") {
			since := time.Now().Add(-time.Hour * 24 * time.Duration(mustGetCmd[int](*cmd, "lookback-days")))
			var err error
			if ret, err = p.Sessions.HistoryEpisodes(&since, args...); err != nil {
				return err
			}
		} else {
			var err error
			if ret, err = p.Sessions.ActiveEpisodes(args...); err != nil {
				return err
			}
		}
		for _, episode := range ret {
			fmt.Println(episode.String())
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getSessionsCmd)
	getSessionsCmd.PersistentFlags().Bool("history", false, "include session history, not just active")
	getSessionsCmd.PersistentFlags().Int("lookback-days", 14, "number of days to look back at viewed history")
}
