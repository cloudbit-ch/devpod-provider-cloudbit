package cmd

import (
	"context"
	"encoding/base64"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/cloudbit"
	"github.com/cloudbit/devpod-provider-cloudbit/pkg/options"
	"github.com/flowswiss/goclient/compute"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		req, err := buildCreateInstanceRequest(options)
		if err != nil {
			return err
		}

		return cloudbit.NewCloudbit(options.Token).CreateInstance(context.Background(), *req)
	},
}

func buildCreateInstanceRequest(opts *options.Options) (*compute.ServerCreate, error) {
	cloudBitClient := cloudbit.NewCloudbit(opts.Token)

	loc, err := cloudBitClient.GetLocationByName(context.Background(), opts.Location)
	if err != nil {
		return nil, err
	}

	image, err := cloudBitClient.GetImageByKey(context.Background(), opts.Image)
	if err != nil {
		return nil, err
	}

	product, err := cloudBitClient.GetProductByName(context.Background(), opts.Product)
	if err != nil {
		return nil, err
	}

	network, err := cloudBitClient.GetNetworkByName(context.Background(), opts.Network)
	if err != nil {
		return nil, err
	}

	keyPair, err := cloudBitClient.CreateKeyPair(context.Background(), opts.MachineID, opts.MachineFolder)
	if err != nil {
		return nil, err
	}

	initScript, err := getInjectKeypairScript(opts.MachineFolder)
	if err != nil {
		return nil, err
	}

	serverCreateReq := &compute.ServerCreate{
		Name:             opts.MachineID,
		LocationID:       loc.ID,
		ImageID:          image.ID,
		ProductID:        product.ID,
		AttachExternalIP: true,
		NetworkID:        network.ID,
		KeyPairID:        keyPair.ID,
		CloudInit:        initScript,
	}
	return serverCreateReq, nil
}

func getInjectKeypairScript(dir string) (string, error) {
	publicKey, err := cloudbit.GetMachinePublicKey(dir)
	if err != nil {
		return "", err
	}

	resultScript := `#!/bin/sh
useradd devpod -d /home/devpod
mkdir -p /home/devpod
if grep -q sudo /etc/groups; then
	usermod -aG sudo devpod
elif grep -q wheel /etc/groups; then
	usermod -aG wheel devpod
fi
echo "devpod ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/91-devpod
mkdir -p /home/devpod/.ssh
echo "` + string(publicKey) + `" >> /home/devpod/.ssh/authorized_keys
chmod 0700 /home/devpod/.ssh
chmod 0600 /home/devpod/.ssh/authorized_keys
chown -R devpod:devpod /home/devpod`

	return base64.StdEncoding.EncodeToString([]byte(resultScript)), nil
}
