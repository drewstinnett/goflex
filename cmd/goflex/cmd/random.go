package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random config-1.yaml [config-2.yaml ...]",
	Short: "Randomize a playlist using the given list of requests",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, flexes, err := loadConfigs(args)
		if err != nil {
			return err
		}

		errs := errgroup.Group{}
		for configIdx, config := range configs {
			slog.Debug("got requests", "requests", config.Randomize)
			if len(config.Randomize) == 0 {
				slog.Warn("config has no randomize requests", "idx", configIdx)
			}
			for showIdx, request := range config.Randomize {
				///idx := idx
				///request := request
				startRandomizer(&errs, showIdx, request, flexes[configIdx])
			}
		}
		return errs.Wait()
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}

func startRandomizer(g *errgroup.Group, idx int, req goflex.RandomizeRequest, f *goflex.Flex) {
	g.Go(func() error {
		logger := slog.With("show-idx", idx)
		logger.Info("randomizing", "playlist", req.Playlist)
		for {
			// TODO: How do we know which plex server to send this to?
			resp, err := f.Playlists.Randomize(req)
			if err != nil {
				return err
			}
			logger.Debug("sleeping for", "duration", resp.SleepFor, "playlist", req.Playlist)
			time.Sleep(resp.SleepFor)
		}
	})
}

func loadConfigs(paths []string) ([]goflex.FlexConfig, []*goflex.Flex, error) {
	cfgs := make([]goflex.FlexConfig, 0, len(paths))
	flexes := make([]*goflex.Flex, 0, len(paths))

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("reading %q: %w", path, err)
		}

		var cfg goflex.FlexConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, nil, fmt.Errorf("parsing %q: %w", path, err)
		}

		flex, err := goflex.New(goflex.WithFlexConfig(cfg))
		if err != nil {
			return nil, nil, fmt.Errorf("initializing flex for %q: %w", path, err)
		}

		cfgs = append(cfgs, cfg)
		flexes = append(flexes, flex)
	}
	return cfgs, flexes, nil
}
