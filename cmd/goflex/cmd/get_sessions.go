package cmd

import (
	"time"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

func stringsToShowTitles(s []string) []goflex.ShowTitle {
	ret := make([]goflex.ShowTitle, len(s))
	for idx, item := range s {
		ret[idx] = goflex.ShowTitle(item)
	}
	return ret
}

// getSessionsCmd represents the random command
var getSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Get session information",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		var ret goflex.EpisodeList
		if mustGetCmd[bool](*cmd, "history") {
			since := time.Now().Add(-time.Hour * 24 * time.Duration(mustGetCmd[int](*cmd, "lookback-days")))
			var err error
			if ret, err = p.Sessions.HistoryEpisodes(&since, stringsToShowTitles(args)...); err != nil {
				return err
			}
		} else {
			var err error
			if ret, err = p.Sessions.ActiveEpisodes(stringsToShowTitles(args)...); err != nil {
				return err
			}
		}
		printEpisodes(ret, mustGetCmd[bool](*cmd, "short"))

		return nil
	},
}

func init() {
	getCmd.AddCommand(getSessionsCmd)
	getSessionsCmd.PersistentFlags().Bool("history", false, "include session history, not just active")
	getSessionsCmd.PersistentFlags().Int("lookback-days", 14, "number of days to look back at viewed history")
	getSessionsCmd.PersistentFlags().BoolP("short", "s", false, "Show short version of the episode (Name S00E00)")
}
