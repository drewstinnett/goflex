package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getNotificationsCmd represents the random command
var getServersCmd = &cobra.Command{
	Use:     "servers",
	Short:   "Get servers",
	Aliases: []string{"server"},
	Args:    cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		got, err := p.Server.Servers()
		if err != nil {
			return err
		}
		gout.MustPrint(got)

		return nil
	},
}

func init() {
	getCmd.AddCommand(getServersCmd)
}
