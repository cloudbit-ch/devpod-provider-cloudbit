package cmd

import "github.com/spf13/cobra"

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve the status of an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
