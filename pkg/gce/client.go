package gce

import (
	"google.golang.org/api/compute/v1"
	"golang.org/x/oauth2/google"
	"context"
)

type GCEClient struct {
	service *compute.Service
}

func NewGCEClient(ctx context.Context) (*GCEClient, error) {
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return &GCEClient{service: service}, nil
}

func (c *GCEClient) StartInstance(projectID, zone, instanceID string) error {
	_, err := c.service.Instances.Start(projectID, zone, instanceID).Context(context.Background()).Do()
	return err
}

func (c *GCEClient) StopInstance(projectID, zone, instanceID string) error {
	_, err := c.service.Instances.Stop(projectID, zone, instanceID).Context(context.Background()).Do()
	return err
}

func (c *GCEClient) ListInstances(projectID, zone string) (*compute.InstanceList, error) {
	return c.service.Instances.List(projectID, zone).Context(context.Background()).Do()
}

