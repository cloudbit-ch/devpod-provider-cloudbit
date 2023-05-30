package cloudbit

import (
	"context"
	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.cloudbit.ch/"
)

type Cloudbit struct {
	computeService compute.ServerService
}

func NewCloudbit(token string) *Cloudbit {
	c := goclient.NewClient(
		goclient.WithBase(baseURL),
		goclient.WithToken(token),
	)
	return &Cloudbit{computeService: compute.NewServerService(c)}
}

func (c *Cloudbit) Init(ctx context.Context) error {
	_, err := c.computeService.List(ctx, goclient.Cursor{})
	if err != nil {
		return errors.Wrap(err, "list compute instances")
	}

	return nil
}

func (c *Cloudbit) Create(ctx context.Context, req compute.ServerCreate) error {
	//c.computeService.Create()
	return nil
}

func (c *Cloudbit) Stop(ctx context.Context, instanceID int) error {
	_, err := c.computeService.Perform(ctx, instanceID, compute.ServerPerform{Action: "stop"})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloudbit) Status(ctx context.Context, instanceID int) (client.Status, error) {
	server, err := c.computeService.Get(ctx, instanceID)
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

func (c *Cloudbit) Delete(ctx context.Context, instanceID int) error {
	return c.computeService.Delete(ctx, instanceID, true)
}
