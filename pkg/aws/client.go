package aws

import (
	"context"
	"sync"

	"github.com/nduyphuong/gorya/pkg/aws/ec2"
	"github.com/nduyphuong/gorya/pkg/aws/options"
)

//go:generate mockery --name Interface --output fake --outpkg fake
type Interface interface {
	EC2() ec2.Interface
}

type client struct {
	ec2 ec2.Interface
}

func New(ctx context.Context, region string, opts ...options.Option) (Interface, error) {
	return getOnce(ctx, region, opts...)
}

var (
	awsClient   *client
	muAwsClient sync.Mutex
)

func getOnce(ctx context.Context, region string, opts ...options.Option) (*client, error) {
	muAwsClient.Lock()
	defer func() {
		muAwsClient.Unlock()
	}()
	if awsClient != nil {
		return awsClient, nil
	}
	var c client
	var err error
	if c.ec2, err = ec2.New(ctx, region, opts...); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) EC2() ec2.Interface { return c.ec2 }
