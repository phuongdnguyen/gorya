package cloudsql_test

import (
	"context"
	"testing"

	"github.com/nduyphuong/gorya/pkg/gcp/cloudsql"
	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/impersonate"
)

type TestData struct {
	TargetPrincipal string
	GCPProjectId    string
	Scopes          []string
}

func TestSmoke(t *testing.T) {
	ctx := context.TODO()
	d := TestData{
		TargetPrincipal: "priv-sa@target-project-397310.iam.gserviceaccount.com",
		GCPProjectId:    "target-project-397310",
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/compute",
		},
	}
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: d.TargetPrincipal,
		Scopes:          d.Scopes,
	})
	assert.NoError(t, err)
	sqlService, err := cloudsql.NewService(ctx, &ts, options.WithProject(d.GCPProjectId))
	assert.NoError(t, err)
	err = sqlService.ChangeStatus(ctx, 0, "foo", "bar")
	assert.NoError(t, err)
}
