package cmd

import (
	"context"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		cloudBitClient := cloudbit.NewCloudbit(options.Token)
		err = cloudBitClient.StartInstanceByName(context.Background(), options.MachineID)
		if err != nil {
			return err
		}

		// wait until running
		for {
			status, err := cloudBitClient.GetStatusByInstanceName(context.Background(), options.MachineID)
			if err != nil {
				log.Default.Errorf("Error retrieving instance status: %v", err)
				break
			}

			if status == client.StatusRunning {
				break
			}

			// make sure we don't spam
			time.Sleep(time.Second)
		}

		return nil
	},
}
