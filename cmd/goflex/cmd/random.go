package cmd

import (
	"log/slog"
	"os"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random REQUEST.yaml",
	Short: "Randomize a playlist using the given list of requests",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := newPlex()

		b, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		var requests goflex.RandomizeRequestList
		if err := yaml.Unmarshal(b, &requests); err != nil {
			return err
		}
		slog.Debug("got requests", "requests", requests)
		for _, request := range requests {
			resp, err := p.Playlists.Randomize(request)
			if err != nil {
				return err
			}
			_ = resp
			gout.MustPrint(resp.SleepFor)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}
