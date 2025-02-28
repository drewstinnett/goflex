package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getNotificationsCmd represents the random command
var getPreferencesCmd = &cobra.Command{
	Use:     "preferences",
	Short:   "Get preferences",
	Aliases: []string{"pref", "prefs", "preference"},
	Args:    cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		got, err := p.Server.Preferences()
		if err != nil {
			return err
		}
		gout.MustPrint(got)

		return nil
	},
}

func init() {
	getCmd.AddCommand(getPreferencesCmd)
}
