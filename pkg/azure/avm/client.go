package avm

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/logging"
	"github.com/nduyphuong/gorya/pkg/azure/options"
)

var ErrInvalidResourceStatus = errors.New("invalid resource status")

type client struct {
	avm  *armcompute.VirtualMachinesClient
	opts options.Options
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
	c.avm = computeClientFactory.NewVirtualMachinesClient()
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
	fmt.Printf("c.opts.TargetResourceGroups: %v\n", c.opts.TargetResourceGroups)
	for _, rg := range c.opts.TargetResourceGroups {
		vmPager := c.avm.NewListPager(rg, &armcompute.VirtualMachinesClientListOptions{})
		for vmPager.More() {
			page, err := vmPager.NextPage(ctx)
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
					if err = c.turnOffInstance(ctx, *item.Name, rg); err != nil {
						logger.Errorf("turn off avm %s in resource group %s %v", *item.Name, rg, err)
						continue
					}
					logger.Infof("turned off avm %s in resource group %s", *item.Name, rg)
				case constants.OnStatus:
					if err = c.turnOnInstance(ctx, *item.Name, rg); err != nil {
						logger.Errorf("turn on avm %s in resource group %s %v", *item.Name, rg, err)
						continue
					}
					logger.Infof("turned on avm %s in resource group %s", *item.Name, rg)
				}
			}
		}
	}
	return nil
}

func (c *client) turnOnInstance(ctx context.Context, name string, resourceGroupName string) error {
	if _, err := c.avm.BeginStart(ctx, resourceGroupName, name, nil); err != nil {
		return err
	}
	return nil
}

func (c *client) turnOffInstance(ctx context.Context, name string, resourceGroupName string) error {
	if _, err := c.avm.BeginPowerOff(ctx, resourceGroupName, name, nil); err != nil {
		return err
	}
	return nil
}
