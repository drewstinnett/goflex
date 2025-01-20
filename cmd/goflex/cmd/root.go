/*
Package cmd is the cmd stuff
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	goflex "github.com/drewstinnett/go-flex"
	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "goflex",
	Short:         "Do better playlist randomization",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Warn("fatal error", "error", err)
		os.Exit(2)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Slog
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          "goflex 🍿 ",
	}
	if verbose {
		opts.Level = log.DebugLevel
	}
	logger := slog.New(log.NewWithOptions(os.Stderr, opts))
	slog.SetDefault(logger)
}

func newPlex() *goflex.Plex {
	opts := []func(*goflex.Plex){
		goflex.WithBaseURL(os.Getenv("PLEX_URL")),
		goflex.WithToken(os.Getenv("PLEX_TOKEN")),
	}
	if os.Getenv("DEBUG_CURL") != "" {
		opts = append(opts, goflex.WithPrintCurl())
	}
	p, err := goflex.New(opts...)
	if err != nil {
		panic(err)
	}
	return p
}
