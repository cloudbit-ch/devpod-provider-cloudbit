package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init called")
		return nil
	},
}
