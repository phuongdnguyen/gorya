package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/nduyphuong/gorya/internal/os"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
*
https://cloudstudio.com.au/2021/05/29/account-structure-aws-vs-azure/
*
*/

type TestData struct {
	TenantId       string
	SubscriptionId string
	TargetVMName   string
	ClientId       string
	ClientSecret   string
	TargetTag      struct {
		Key, Value string
	}
}

func TestNewPool(t *testing.T) {
	ctx := context.TODO()
	d := TestData{
		TenantId:       os.MustGetEnv("AZURE_TENANT_ID"),
		SubscriptionId: os.MustGetEnv("AZURE_TARGET_SUBSCRIPTION_ID"),
		ClientId:       os.MustGetEnv("AZURE_CLIENT_ID"),
		ClientSecret:   os.MustGetEnv("AZURE_CLIENT_SECRET"),
		TargetVMName:   "test-vm",
		TargetTag:      struct{ Key, Value string }{Key: "foo", Value: "bar"},
	}

	conn, err := azidentity.NewDefaultAzureCredential(nil)
	armSubscription, err := armsubscriptions.NewClient(conn, nil)
	assert.NoError(t, err)
	subPager := armSubscription.NewListPager(nil)
	credentialRef := map[string]bool{}
	for subPager.More() {
		page, err := subPager.NextPage(ctx)
		assert.NoError(t, err)
		for _, subscription := range page.Value {
			credentialRef[*subscription.SubscriptionID] = true
		}
	}
	azurePool, err := NewPool(ctx, credentialRef)
	assert.NoError(t, err)
	azClient, exist := azurePool.GetForCredential(d.SubscriptionId)
	assert.True(t, exist)
	err = azClient.AVM().ChangeStatus(ctx, 0, d.TargetTag.Key, d.TargetTag.Value)
	assert.NoError(t, err)
}
