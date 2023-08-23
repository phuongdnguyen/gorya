package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nduyphuong/gorya/pkg/aws/ec2"
	"github.com/nduyphuong/gorya/pkg/aws/options"
)

//go:generate mockery --name Interface
type Interface interface {
	EC2() ec2.Interface
}

type client struct {
	ec2  ec2.Interface
	opts options.Options
}

func New(ctx context.Context, opts ...options.Option) (Interface, error) {
	return getOnce(ctx, opts...)
}

var (
	awsClient   *client
	muAwsClient sync.Mutex
)

func getOnce(ctx context.Context, opts ...options.Option) (*client, error) {
	muAwsClient.Lock()
	defer func() {
		muAwsClient.Unlock()
	}()
	if awsClient != nil {
		return awsClient, nil
	}
	var c client
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	awsEndpoint := c.opts.AwsEndpoint
	awsRegion := c.opts.AwsRegion
	// custom resolver so we can testing locally with localstack
	customResolverWithOptions := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if awsEndpoint != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           awsEndpoint,
					SigningRegion: awsRegion,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.opts.AwsEndpoint),
		config.WithEndpointResolverWithOptions(customResolverWithOptions),
	)
	if err != nil {
		return nil, err
	}
	if c.ec2, err = ec2.NewFromConfig(cfg); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) EC2() ec2.Interface { return c.ec2 }
