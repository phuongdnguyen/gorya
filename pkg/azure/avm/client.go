package avm

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/logging"
	"github.com/nduyphuong/gorya/pkg/azure/options"
)

var ErrInvalidResourceStatus = errors.New("invalid resource status")

type client struct {
	avm              *armcompute.VirtualMachinesClient
	armResourceGraph *armresourcegraph.Client
	opts             options.Options
}
type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

func New(conn azcore.TokenCredential, opts ...options.Option) (Interface, error) {
	c := &client{}
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	computeClientFactory, err := armcompute.NewClientFactory(c.opts.SubscriptionId, conn, nil)
	resourceGraphClientFactory, err := armresourcegraph.NewClientFactory(conn, nil)
	c.avm = computeClientFactory.NewVirtualMachinesClient()
	c.armResourceGraph = resourceGraphClientFactory.NewClient()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error) {
	logger := logging.LoggerFromContext(ctx)
	if to != constants.OffStatus && to != constants.OnStatus {
		return ErrInvalidResourceStatus
	}
	pager := c.avm.NewListPager(c.opts.ResourceGroupName, &armcompute.VirtualMachinesClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, item := range page.Value {
			if item.Tags == nil {
				continue
			}
			if *item.Tags[tagKey] != tagValue {
				continue
			}
			switch to {
			case constants.OffStatus:
				if err = c.turnOffInstance(ctx, *item.Name); err != nil {
					logger.Errorf("turn off avm %s %v", *item.Name, err)
				}
			case constants.OnStatus:
				if err = c.turnOnInstance(ctx, *item.Name); err != nil {
					logger.Errorf("turn on avm %s %v", *item.Name, err)
				}
			}
		}
	}
	return nil
}

func (c *client) turnOnInstance(ctx context.Context, name string) error {
	if _, err := c.avm.BeginStart(ctx, c.opts.ResourceGroupName, name, nil); err != nil {
		return err
	}
	return nil
}

func (c *client) turnOffInstance(ctx context.Context, name string) error {
	if _, err := c.avm.BeginPowerOff(ctx, c.opts.ResourceGroupName, name, nil); err != nil {
		return err
	}
	return nil
}
