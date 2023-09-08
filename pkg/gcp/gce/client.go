package gce

import (
	"context"
	"errors"

	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"github.com/nduyphuong/gorya/pkg/gcp/utils"
	"golang.org/x/oauth2"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

//go:generate mockery --name Interface
type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

type client struct {
	gce  *compute.Service
	opts options.Options
}

func NewService(ctx context.Context, ts *oauth2.TokenSource, opts ...options.Option) (*client, error) {
	var err error
	c := &client{}
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	c.gce, err = compute.NewService(ctx, option.WithTokenSource(*ts))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) error {
	if to != 0 && to != 1 {
		return errors.New("to must have value of 0 or 1")
	}
	var err error
	zoneListResp, err := c.gce.Zones.List(c.opts.Project).Do()
	if err != nil {
		return err
	}
	var zones []string
	for _, zone := range zoneListResp.Items {
		zones = append(zones, zone.Description)
	}
	tagFilter := utils.GetFilter(tagKey, tagValue)
	for _, zone := range zones {
		instanceListResp, err := c.gce.Instances.List(c.opts.Project, zone).Context(ctx).Filter(tagFilter).Do()
		if err != nil {
			return err
		}
		for _, instance := range instanceListResp.Items {
			switch to {
			case 0:
				c.gce.Instances.Stop(c.opts.Project, zone, instance.Name).Do()
			case 1:
				c.gce.Instances.Start(c.opts.Project, zone, instance.Name).Do()
			}
		}
	}
	return nil
}
