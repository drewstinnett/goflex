package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getLibraryCmd represents the random command
var getTokenCmd = &cobra.Command{
	Use:   "token USERNAME PASSWORD",
	Short: "Get a new token from username and password",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		p := newPlex()

		got, err := p.Authentication.Token(args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Println(got)
		return nil
	},
}

func init() {
	getCmd.AddCommand(getTokenCmd)
}
