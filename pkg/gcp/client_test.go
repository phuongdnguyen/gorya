package gcp

import (
	"context"
	"testing"

	"github.com/nduyphuong/gorya/internal/os"
	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx := context.TODO()
	project := os.GetEnv("PROJECT", "target-project-397310")
	c, err := new(ctx,
		options.WithImpersonatedServiceAccountEmail("priv-sa@target-project-397310.iam.gserviceaccount.com"),
		options.WithProject(project),
	)
	assert.NoError(t, err)
	err = c.GCE().ChangeStatus(ctx, 0, "foo", "bar")
	assert.NoError(t, err)
}
