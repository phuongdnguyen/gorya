package aws

import (
	"context"
	"testing"

	"github.com/nduyphuong/gorya/internal/os"
	awsOptions "github.com/nduyphuong/gorya/pkg/aws/options"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx := context.TODO()
	awsRegion := os.GetEnv("AWS_REGION", "ap-southeast-1")
	awsEndpoint := os.GetEnv("AWS_ENDPOINT", "")
	c, err := New(ctx,
		awsOptions.WithRegion(awsRegion),
		awsOptions.WithEndpoint(awsEndpoint),
	)
	assert.NoError(t, err)
	err = c.EC2().ChangeStatus(ctx, 0, "foo", "bar")
	assert.NoError(t, err)
}
