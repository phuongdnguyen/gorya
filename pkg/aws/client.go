package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nduyphuong/gorya/internal/constants"
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

type ClientPool struct {
	credToClient map[string]Interface
}

var (
	lock sync.Mutex
)

func NewPool(ctx context.Context, credentialRefs map[string]bool,
	opts ...options.Option) (*ClientPool,
	error) {
	lock.Lock()
	defer lock.Unlock()
	b := &ClientPool{
		credToClient: make(map[string]Interface),
	}
	for cred := range credentialRefs {
		if _, ok := b.credToClient[cred]; !ok {
			c, err := new(ctx, append(opts, options.WithRoleArn(cred))...)
			if err != nil {
				return nil, err
			}
			b.credToClient[cred] = c
		}
	}
	return b, nil
}

func (b *ClientPool) GetForCredential(name string) (Interface, bool) {
	if name == constants.Default {
		return b.credToClient[constants.Default], true
	}
	i, ok := b.credToClient[name]
	if !ok {
		return nil, false
	}
	fmt.Printf("got client from pool for %s", name)
	return i, true
}

func new(ctx context.Context, opts ...options.Option) (*client, error) {
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
	if c.opts.AwsRoleArn != constants.Default {
		stsClient := sts.NewFromConfig(cfg)
		provider := stscreds.NewAssumeRoleProvider(stsClient, c.opts.AwsRoleArn)
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}
	if c.ec2, err = ec2.NewFromConfig(cfg); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) EC2() ec2.Interface { return c.ec2 }
