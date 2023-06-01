package cmd

import (
	"context"
	"fmt"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"
	"github.com/spf13/cobra"
	"os"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve the status of an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		status, err := cloudbit.NewCloudbit(options.Token).GetStatusByInstanceName(context.Background(), options.MachineID)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, status)
		return err
	},
}
