package ec2

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nduyphuong/gorya/pkg/aws/options"
)

//go:generate mockery --name Interface --output fake --outpkg fake
type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

type Client struct {
	ec2  *ec2.Client
	opts options.Options
}

func New(ctx context.Context, region string, opts ...options.Option) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	c := &Client{}
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	c.ec2 = ec2.NewFromConfig(cfg)
	return c, nil
}

func (c *Client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error) {
	print(to)
	if to != 0 && to != 1 {
		return errors.New("to must have value of 0 or 1")
	}
	var describeInstancesOutput *ec2.DescribeInstancesOutput
	filterName := fmt.Sprintf("tag:%s", tagKey)
	describeInstancesOutput, err = c.ec2.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   &filterName,
				Values: []string{tagValue},
			},
		},
	})
	if err != nil {
		return err
	}
	var instancesIds []string
	if describeInstancesOutput != nil {
		for _, r := range describeInstancesOutput.Reservations {
			for _, i := range r.Instances {
				instancesIds = append(instancesIds, *i.InstanceId)
			}
		}
	}
	switch to {
	case 0:
		if _, err = c.ec2.StopInstances(ctx, &ec2.StopInstancesInput{
			InstanceIds: instancesIds,
		}); err != nil {
			return err
		}
	case 1:
		if _, err = c.ec2.StartInstances(ctx, &ec2.StartInstancesInput{
			InstanceIds: instancesIds,
		}); err != nil {
			return err
		}
	}
	return nil
}
