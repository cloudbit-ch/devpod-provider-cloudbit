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

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		cloudBitClient := cloudbit.NewCloudbit(options.Token)
		err = cloudBitClient.Stop(context.Background(), options.MachineID)
		if err != nil {
			return err
		}

		// wait until stopped
		for {
			status, err := cloudBitClient.Status(context.Background(), options.MachineID)
			if err != nil {
				log.Default.Errorf("Error retrieving instance status: %v", err)
				break
			}

			if status == client.StatusStopped {
				break
			}

			// make sure we don't spam
			time.Sleep(time.Second)
		}

		return nil
	},
}
