package cmd

import (
	"os"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "devpod-provider-cloudbit",
	Short:         "Cloudbit Provider commands",
	SilenceErrors: true,
	SilenceUsage:  true,

	PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
		log.Default.MakeRaw()
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commandCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(stopCmd)
}
