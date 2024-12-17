package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/LukeHagar/plexgo"
	plexrando "github.com/drewstinnett/plex-truerandom"
	"github.com/spf13/cobra"
)

// historyCmd represents the random command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Session history",
	RunE: func(cmd *cobra.Command, _ []string) error {
		p := plexrando.New(plexrando.WithAPI(plexgo.New(
			plexgo.WithSecurity(os.Getenv("PLEX_TOKEN")),
			plexgo.WithServerURL(os.Getenv("PLEX_URL")),
			plexgo.WithClientID("313FF6D7-5795-45E3-874F-B8FCBFD5E587"),
			plexgo.WithClientName("plex-trueget"),
			plexgo.WithClientVersion("0.0.1"),
		)))

		viewed, err := p.Viewed(mustGetCmd[string](*cmd, "library"), time.Now().Add(-time.Hour*24*14))
		if err != nil {
			return err
		}

		for _, episode := range viewed {
			fmt.Println(episode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
}
