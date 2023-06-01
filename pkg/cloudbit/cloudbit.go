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

func (c *Cloudbit) Create(ctx context.Context, req compute.ServerCreate) error {
	_, err := c.computeService.Create(ctx, req)
	return err
}

func (c *Cloudbit) Stop(ctx context.Context, machineID string) error {
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

func (c *Cloudbit) Start(ctx context.Context, machineID string) error {
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

func (c *Cloudbit) Status(ctx context.Context, machineID string) (client.Status, error) {
	server, err := c.GetInstanceByName(ctx, machineID)
	if err != nil {
		return "", err
	}

	server, err = c.computeService.Get(ctx, server.ID)
	if err != nil {
		return client.StatusNotFound, err
	}

	switch server.Status.Name {
	case client.StatusRunning:
		return client.StatusRunning, nil
	case client.StatusStopped:
		return client.StatusStopped, nil
	case client.StatusBusy:
		return client.StatusBusy, nil
	}

	return client.StatusNotFound, nil
}

func (c *Cloudbit) Delete(ctx context.Context, machineID string) error {
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

func (c *Cloudbit) GetInstanceByName(ctx context.Context, machineID string) (compute.Server, error) {
	serverList, err := c.computeService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return compute.Server{}, err
	}

	for _, server := range serverList.Items {
		if server.Name == machineID {
			return server, nil
		}
	}

	return compute.Server{}, errors.New("instance name not found")
}

func (c *Cloudbit) GetInstancePublicIPByName(ctx context.Context, machineID string) (string, error) {
	elasticIPList, err := c.elasticIPService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return "", err
	}

	for _, elasticIP := range elasticIPList.Items {
		if elasticIP.Attachment.Name == machineID {
			return elasticIP.PublicIP, nil
		}
	}

	return "", errors.New("instance public ip not found")
}

func (c *Cloudbit) GetLocationByName(ctx context.Context, name string) (common.Location, error) {
	locationList, err := c.locationService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return common.Location{}, err
	}

	for _, location := range locationList.Items {
		if location.Name == name {
			return location, nil
		}
	}

	return common.Location{}, errors.New("compute location not found")
}

func (c *Cloudbit) GetImageByKey(ctx context.Context, key string) (compute.Image, error) {
	imageList, err := c.imageService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return compute.Image{}, err
	}

	for _, image := range imageList.Items {
		if image.Key == key {
			return image, nil
		}
	}

	return compute.Image{}, errors.New("compute image not found")
}

func (c *Cloudbit) GetProductByName(ctx context.Context, name string) (common.Product, error) {
	productList, err := c.productService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return common.Product{}, err
	}

	for _, product := range productList.Items {
		if product.Name == name {
			return product, nil
		}
	}

	return common.Product{}, errors.New("compute product not found")
}

func (c *Cloudbit) GetNetworkByName(ctx context.Context, name string) (compute.Network, error) {
	networkList, err := c.networkService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return compute.Network{}, err
	}

	for _, network := range networkList.Items {
		if network.Name == name {
			return network, nil
		}
	}

	return compute.Network{}, errors.New("compute network not found")
}

func (c *Cloudbit) GetKeyPairByName(ctx context.Context, name string) (compute.KeyPair, error) {
	keyPairList, err := c.keyPairService.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return compute.KeyPair{}, err
	}

	for _, keyPair := range keyPairList.Items {
		if keyPair.Name == name {
			return keyPair, nil
		}
	}

	return compute.KeyPair{}, errors.New("compute keypair not found")
}

func (c *Cloudbit) CreateKeyPair(ctx context.Context, name string, dir string) (compute.KeyPair, error) {
	publicKey, err := GetMachinePublicKey(dir)
	if err != nil {
		return compute.KeyPair{}, err
	}

	// reuse the keypair if already exists
	keyPair, err := c.GetKeyPairByName(ctx, name)
	if err != nil {
		keyPair, err = c.keyPairService.Create(ctx, compute.KeyPairCreate{
			Name:      name,
			PublicKey: publicKey,
		})
		if err != nil {
			return compute.KeyPair{}, err
		}

		return keyPair, nil
	}

	return keyPair, nil
}

func (c *Cloudbit) DeleteKeyPairByName(ctx context.Context, name string) error {
	keyPair, err := c.GetKeyPairByName(ctx, name)
	if err != nil {
		return err
	}

	return c.keyPairService.Delete(ctx, keyPair.ID)
}

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
