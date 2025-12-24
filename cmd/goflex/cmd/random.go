package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

		// Set up context with signal handling for graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigChan
			slog.Info("received signal, initiating graceful shutdown", "signal", sig)
			cancel()
		}()

		errs, ctx := errgroup.WithContext(ctx)
		for configIdx, config := range configs {
			slog.Debug("got requests", "requests", config.Randomize)
			if len(config.Randomize) == 0 {
				slog.Warn("config has no randomize requests", "idx", configIdx)
			}
			for showIdx, request := range config.Randomize {
				startRandomizer(ctx, errs, showIdx, request, flexes[configIdx])
			}
		}
		err = errs.Wait()
		if err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		slog.Info("shutdown complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}

const (
	maxRetries       = 10
	initialBackoff   = 30 * time.Second
	maxBackoff       = 30 * time.Minute
	backoffMultiplier = 2.0
)

// isFatalError returns true if the error should not be retried.
// Fatal errors include configuration problems and missing resources.
func isFatalError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// These indicate configuration or setup problems that won't resolve with retries
	fatalPatterns := []string{
		"show does not exist",
		"playlist must not be empty",
		"series muset not be empty", // typo preserved from original code
		"must set plex baseurl",
		"must set token",
	}
	for _, pattern := range fatalPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}

func calculateBackoff(consecutiveErrors int) time.Duration {
	if consecutiveErrors <= 0 {
		return initialBackoff
	}
	backoff := float64(initialBackoff)
	for i := 0; i < consecutiveErrors-1; i++ {
		backoff *= backoffMultiplier
		if backoff > float64(maxBackoff) {
			return maxBackoff
		}
	}
	return time.Duration(backoff)
}

func startRandomizer(ctx context.Context, g *errgroup.Group, idx int, req goflex.RandomizeRequest, f *goflex.Flex) {
	g.Go(func() error {
		logger := slog.With("show-idx", idx, "playlist", req.Playlist)
		logger.Info("starting randomizer")

		var consecutiveErrors int

		for {
			select {
			case <-ctx.Done():
				logger.Info("randomizer stopping due to context cancellation")
				return ctx.Err()
			default:
			}

			resp, err := f.Playlists.Randomize(req)
			if err != nil {
				// Check if this is a fatal error that shouldn't be retried
				if isFatalError(err) {
					logger.Error("fatal error, stopping randomizer", "error", err)
					return fmt.Errorf("fatal error for playlist %q: %w", req.Playlist, err)
				}

				consecutiveErrors++
				backoff := calculateBackoff(consecutiveErrors)

				if consecutiveErrors >= maxRetries {
					logger.Error("max retries exceeded, stopping randomizer",
						"error", err,
						"consecutive_errors", consecutiveErrors)
					return fmt.Errorf("max retries (%d) exceeded for playlist %q: %w", maxRetries, req.Playlist, err)
				}

				logger.Warn("randomize failed, will retry",
					"error", err,
					"consecutive_errors", consecutiveErrors,
					"backoff", backoff)

				select {
				case <-ctx.Done():
					logger.Info("randomizer stopping during backoff")
					return ctx.Err()
				case <-time.After(backoff):
					continue
				}
			}

			// Success - reset error counter
			if consecutiveErrors > 0 {
				logger.Info("recovered after errors", "previous_consecutive_errors", consecutiveErrors)
			}
			consecutiveErrors = 0

			logger.Debug("sleeping until next check", "duration", resp.SleepFor)

			select {
			case <-ctx.Done():
				logger.Info("randomizer stopping during sleep")
				return ctx.Err()
			case <-time.After(resp.SleepFor):
			}
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
