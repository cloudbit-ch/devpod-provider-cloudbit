package cloudbit

import (
	"context"
	"encoding/base64"
	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/common"
	"github.com/flowswiss/goclient/compute"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.cloudbit.ch/"
)

type Cloudbit struct {
	computeService   compute.ServerService
	elasticIPService compute.ElasticIPService
	imageService     compute.ImageService
	keyPairService   compute.KeyPairService
	networkService   compute.NetworkService
	locationService  common.LocationService
	productService   common.ProductService
}

func NewCloudbit(token string) *Cloudbit {
	c := goclient.NewClient(
		goclient.WithBase(baseURL),
		goclient.WithToken(token),
	)
	return &Cloudbit{
		computeService:   compute.NewServerService(c),
		elasticIPService: compute.NewElasticIPService(c),
		imageService:     compute.NewImageService(c),
		keyPairService:   compute.NewKeyPairService(c),
		networkService:   compute.NewNetworkService(c),
		locationService:  common.NewLocationService(c),
		productService:   common.NewProductService(c),
	}
}

func (c *Cloudbit) Init(ctx context.Context) error {
	_, err := c.computeService.List(ctx, goclient.Cursor{})
	if err != nil {
		return errors.Wrap(err, "list compute instances")
	}

	return nil
}

// CreateInstance created the specified devpod machine instance
func (c *Cloudbit) CreateInstance(ctx context.Context, req compute.ServerCreate) error {
	// re-use instance if exists
	_, err := c.GetInstanceByName(ctx, req.Name)
	if err == nil {
		return nil
	}

	_, err = c.computeService.Create(ctx, req)
	return err
}

// StartInstanceByName starts the specified devpod machine instance
func (c *Cloudbit) StartInstanceByName(ctx context.Context, machineID string) error {
	server, err := c.GetInstanceByName(ctx, machineID)
	if err != nil {
		return err
	}

	_, err = c.computeService.Perform(ctx, server.ID, compute.ServerPerform{Action: "start"})
	if err != nil {
		return err
	}

	return nil
}

// StopInstanceByName stops the specified devpod machine instance
func (c *Cloudbit) StopInstanceByName(ctx context.Context, machineID string) error {
	server, err := c.GetInstanceByName(ctx, machineID)
	if err != nil {
		return err
	}

	_, err = c.computeService.Perform(ctx, server.ID, compute.ServerPerform{Action: "stop"})
	if err != nil {
		return err
	}

	return nil
}

// DeleteInstanceByName deleted the specified devpod machine instance
func (c *Cloudbit) DeleteInstanceByName(ctx context.Context, machineID string) error {
	server, err := c.GetInstanceByName(ctx, machineID)
	if err != nil {
		return err
	}

	// remove instance keypair
	err = c.DeleteKeyPairByName(ctx, machineID)
	if err != nil {
		return err
	}

	return c.computeService.Delete(ctx, server.ID, true)
}

// GetStatusByInstanceName retrieves the status of the specified devpod machine instance
func (c *Cloudbit) GetStatusByInstanceName(ctx context.Context, machineID string) (client.Status, error) {
	server, err := c.GetInstanceByName(ctx, machineID)
	if err != nil {
		return "", err
	}

	instance, err := c.computeService.Get(ctx, server.ID)
	if err != nil {
		return client.StatusNotFound, err
	}

	switch instance.Status.Name {
	case client.StatusRunning:
		return client.StatusRunning, nil
	case client.StatusStopped:
		return client.StatusStopped, nil
	case client.StatusBusy:
		return client.StatusBusy, nil
	}

	return client.StatusNotFound, nil
}

// GetInstanceByName retrieves the devpod instance with the specified name
func (c *Cloudbit) GetInstanceByName(ctx context.Context, machineID string) (*compute.Server, error) {
	serverList, err := c.computeService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, server := range serverList.Items {
		if server.Name == machineID {
			return &server, nil
		}
	}

	return nil, errors.New("instance name not found")
}

// GetElasticIPByInstanceName retrieves the Elastic IP associated with the specified instance name
func (c *Cloudbit) GetElasticIPByInstanceName(ctx context.Context, machineID string) (*compute.ElasticIP, error) {
	elasticIPList, err := c.elasticIPService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, elasticIP := range elasticIPList.Items {
		if elasticIP.Attachment.Name == machineID {
			return &elasticIP, nil
		}
	}

	return nil, errors.New("instance public ip not found")
}

// GetLocationByName retrieves the compute location with the specified location name
func (c *Cloudbit) GetLocationByName(ctx context.Context, name string) (*common.Location, error) {
	locationList, err := c.locationService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, location := range locationList.Items {
		if location.Name == name {
			return &location, nil
		}
	}

	return nil, errors.New("compute location not found")
}

// GetImageByKey retrieves the compute images with the specified image key
func (c *Cloudbit) GetImageByKey(ctx context.Context, key string) (*compute.Image, error) {
	imageList, err := c.imageService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, image := range imageList.Items {
		if image.Key == key {
			return &image, nil
		}
	}

	return nil, errors.New("compute image not found")
}

// GetProductByName retrieves the compute product with the specified product name
func (c *Cloudbit) GetProductByName(ctx context.Context, name string) (*common.Product, error) {
	productList, err := c.productService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, product := range productList.Items {
		if product.Name == name {
			return &product, nil
		}
	}

	return nil, errors.New("compute product not found")
}

// GetNetworkByName retrieves the compute network with the specified network name
func (c *Cloudbit) GetNetworkByName(ctx context.Context, name string) (*compute.Network, error) {
	networkList, err := c.networkService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, network := range networkList.Items {
		if network.Name == name {
			return &network, nil
		}
	}

	return nil, errors.New("compute network not found")
}

// CreateKeyPair creates a new compute key pair
func (c *Cloudbit) CreateKeyPair(ctx context.Context, name string, dir string) (*compute.KeyPair, error) {
	publicKey, err := GetMachinePublicKey(dir)
	if err != nil {
		return nil, err
	}

	// reuse the keypair if already exists
	keyPair, err := c.GetKeyPairByName(ctx, name)
	if err != nil {
		newKeyPair, err := c.keyPairService.Create(ctx, compute.KeyPairCreate{
			Name:      name,
			PublicKey: publicKey,
		})
		if err != nil {
			return nil, err
		}

		return &newKeyPair, nil
	}

	return keyPair, nil
}

// GetKeyPairByName retrieves the compute key pair with the specified pair name
func (c *Cloudbit) GetKeyPairByName(ctx context.Context, name string) (*compute.KeyPair, error) {
	keyPairList, err := c.keyPairService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, keyPair := range keyPairList.Items {
		if keyPair.Name == name {
			return &keyPair, nil
		}
	}

	return nil, errors.New("compute keypair not found")
}

// DeleteKeyPairByName delete the compute key pair with the specified pair name
func (c *Cloudbit) DeleteKeyPairByName(ctx context.Context, name string) error {
	keyPair, err := c.GetKeyPairByName(ctx, name)
	if err != nil {
		return err
	}

	return c.keyPairService.Delete(ctx, keyPair.ID)
}

// GetMachinePublicKey retrieves the machine's public key from the specified directory
func GetMachinePublicKey(dir string) (string, error) {
	publicKeyBase, err := ssh.GetPublicKeyBase(dir)
	if err != nil {
		return "", err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
}
