package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

// getNotificationsCmd represents the random command
var getAccountsCmd = &cobra.Command{
	Use:     "accounts",
	Short:   "Get accounts",
	Aliases: []string{"account"},
	Args:    cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		p := newPlex()

		got, err := p.Server.Accounts()
		if err != nil {
			return err
		}
		return gout.Print(got)
	},
}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
