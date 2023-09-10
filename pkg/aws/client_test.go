package aws

import (
	"context"
	"testing"

	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/os"
	awsOptions "github.com/nduyphuong/gorya/pkg/aws/options"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx := context.TODO()
	awsRegion := os.GetEnv(constants.ENV_AWS_REGION, "ap-southeast-1")
	awsEndpoint := os.GetEnv(constants.ENV_AWS_ENDPOINT, "")
	c, err := new(ctx,
		awsOptions.WithRegion(awsRegion),
		awsOptions.WithEndpoint(awsEndpoint),
	)
	assert.NoError(t, err)
	err = c.EC2().ChangeStatus(ctx, 0, "foo", "bar")
	assert.NoError(t, err)
}
