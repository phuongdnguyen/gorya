package ec2

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
)

func TestClient_ChangeStatus(t *testing.T) {
	var err error
	awsRegion := "ap-southeast-1"
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	assert.NoError(t, err)
	client, err := NewFromConfig(cfg)
	assert.NoError(t, err)
	err = client.ChangeStatus(ctx, 1, "phuong", "test")
	assert.NoError(t, err)

}
