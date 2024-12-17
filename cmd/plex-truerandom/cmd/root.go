/*
Package cmd is the cmd stuff
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/LukeHagar/plexgo"
	"github.com/charmbracelet/log"
	plexrando "github.com/drewstinnett/plex-truerandom"
	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "plex-truerandom",
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Slog
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          "plex-truerandom 🍿 ",
	}

	// Zerolog deprecated
	if verbose {
		opts.Level = log.DebugLevel
	}
	logger := slog.New(log.NewWithOptions(os.Stderr, opts))
	slog.SetDefault(logger)
}

func newPlex() *plexrando.Plex {
	return plexrando.New(plexrando.WithAPI(plexgo.New(
		plexgo.WithSecurity(os.Getenv("PLEX_TOKEN")),
		plexgo.WithServerURL(os.Getenv("PLEX_URL")),
		plexgo.WithClientID("313FF6D7-5795-45E3-874F-B8FCBFD5E587"),
		plexgo.WithClientName("plex-trueget"),
		plexgo.WithClientVersion("0.0.1"),
	)))
}
