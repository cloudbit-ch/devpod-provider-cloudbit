package cmd

import (
	"context"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		req, err := buildCreateInstanceRequest(options)
		if err != nil {
			return err
		}

		return cloudbit.NewCloudbit(options.Token).Create(context.Background(), req)
	},
}
