package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getNotificationsCmd represents the random command
var getCapabilitiesCmd = &cobra.Command{
	Use:     "capabilities",
	Short:   "Get capabilities",
	Aliases: []string{"cap", "capability"},
	Args:    cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		got, err := p.Server.Capabilities()
		if err != nil {
			return err
		}
		gout.MustPrint(got)

		return nil
	},
}

func init() {
	getCmd.AddCommand(getCapabilitiesCmd)
}
