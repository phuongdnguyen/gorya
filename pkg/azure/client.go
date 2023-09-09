package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nduyphuong/gorya/pkg/azure/avm"
	"github.com/nduyphuong/gorya/pkg/azure/options"
	"sync"
)

type Interface interface {
	AVM() avm.Interface
}

type ClientPool struct {
	credToClient map[string]Interface
}

type client struct {
	avm  avm.Interface
	opts options.Options
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
			c, err := new(append(opts, options.WithTenantId(cred))...)
			if err != nil {
				return nil, err
			}
			b.credToClient[cred] = c
		}
	}
	return b, nil
}

func new(opts ...options.Option) (*client, error) {
	var c client
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	conn, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		//load from env
		// AZURE_TENANT_ID: ID of the service principal's tenant. Also called its "directory" ID.
		//
		// AZURE_CLIENT_ID: the service principal's client ID
		//
		// AZURE_CLIENT_SECRET: one of the service principal's client secrets
		//target tenant
		TenantID: c.opts.TenantId,
	})
	if err != nil {
		return nil, err
	}
	c.avm, err = avm.New(conn, opts...)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) AVM() avm.Interface { return c.avm }
