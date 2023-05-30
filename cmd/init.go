package cmd

import (
	"context"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := options.FromEnv(true)
		if err != nil {
			return err
		}

		return cloudbit.NewCloudbit(options.Token).Init(context.Background())
	},
}
