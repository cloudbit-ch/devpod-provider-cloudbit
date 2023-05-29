package cmd

import "github.com/spf13/cobra"

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
