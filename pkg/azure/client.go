package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nduyphuong/gorya/internal/constants"
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

// NewPool return a pool of client identified by subscriptionId
func NewPool(ctx context.Context, credentialRefs map[string]bool,
	opts ...options.Option) (*ClientPool,
	error) {
	lock.Lock()
	defer lock.Unlock()
	b := &ClientPool{
		credToClient: make(map[string]Interface),
	}
	conn, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	for subscription := range credentialRefs {
		//cred = subscription
		if _, ok := b.credToClient[subscription]; !ok {
			resourceGroupClient, err := armresources.NewResourceGroupsClient(subscription, conn, nil)
			if err != nil {
				return nil, err
			}
			rPager := resourceGroupClient.NewListPager(nil)
			var resourceGroups []string
			for rPager.More() {
				page, err := rPager.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, resourceGroup := range page.Value {
					resourceGroups = append(resourceGroups, *resourceGroup.Name)
				}
			}
			c, err := new(conn, append(opts, options.WithSubscriptionId(subscription), options.WithTargetResourceGroups(resourceGroups))...)
			if err != nil {
				return nil, err
			}
			b.credToClient[subscription] = c
		}
	}
	return b, nil
}

// new return a client for a subscriptionId
func new(conn azcore.TokenCredential, opts ...options.Option) (*client, error) {
	avm, err := avm.New(conn, opts...)
	if err != nil {
		return nil, err
	}
	c := client{
		avm: avm,
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (b *ClientPool) GetForCredential(name string) (Interface, bool) {
	if name == constants.Default {
		panic("az subscription id must not be empty")
	}
	i, ok := b.credToClient[name]
	if !ok {
		return nil, false
	}
	fmt.Printf("got client from pool for %s\n", name)
	return i, true
}

func (c *client) AVM() avm.Interface { return c.avm }
