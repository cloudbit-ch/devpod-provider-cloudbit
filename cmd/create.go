package cmd

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
