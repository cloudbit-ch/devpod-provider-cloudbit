package cmd

import "github.com/spf13/cobra"

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Run a command on the instance",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
