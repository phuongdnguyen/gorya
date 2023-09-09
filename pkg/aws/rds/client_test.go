package rds_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nduyphuong/gorya/pkg/aws/rds"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-southeast-1"),
	)
	assert.NoError(t, err)
	c := rds.NewFromConfig(cfg)
	err = c.ChangeStatus(ctx, 1, "foo", "bar")
	assert.NoError(t, err)
}
