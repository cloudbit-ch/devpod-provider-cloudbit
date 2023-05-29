package cmd

import "github.com/spf13/cobra"

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
