package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getLibraryCmd represents the random command
var getLibraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Get library",
	Args:  cobra.ExactArgs(0),
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		items, err := p.Libraries()
		if err != nil {
			return err
		}
		gout.MustPrint(items)
		return nil
	},
}

func init() {
	getCmd.AddCommand(getLibraryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getLibraryCmd.PersistentFlags().String("library", "TV Shows", "Library of the TV Show we are randomizing")
	// getLibraryCmd.PersistentFlags().String("title", "American Dad!", "Name of the show to include in this playlist")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getLibraryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
