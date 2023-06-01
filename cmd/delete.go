package cmd

import (
	"context"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		return cloudbit.NewCloudbit(options.Token).Delete(context.Background(), options.MachineID)
	},
}
