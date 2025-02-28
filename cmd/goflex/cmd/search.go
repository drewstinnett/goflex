package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getNotificationsCmd represents the random command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search libraries for something",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		got, err := p.Server.Search(args[0])
		if err != nil {
			return err
		}
		episodes, err := got.Episodes()
		if err != nil {
			return err
		}
		fmt.Println("Episodes:")
		for _, episode := range episodes {
			fmt.Println(episode.String())
		}

		fmt.Println("Shows:")
		shows, err := got.Shows()
		if err != nil {
			return err
		}
		for _, show := range shows {
			fmt.Println(show.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
