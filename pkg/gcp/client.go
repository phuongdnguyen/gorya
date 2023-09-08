package gcp

import (
	"context"
	"fmt"
	"sync"

	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/pkg/gcp/cloudsql"
	"github.com/nduyphuong/gorya/pkg/gcp/gce"
	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"google.golang.org/api/impersonate"
)

//go:generate mockery --name Interface
type Interface interface {
	GCE() gce.Interface
	CloudSQL() cloudsql.Interface
}

type client struct {
	gce      gce.Interface
	cloudSql cloudsql.Interface
	opts     options.Options
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
			c, err := new(ctx, append(opts, options.WithImpersonatedServiceAccountEmail(cred))...)
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
	fmt.Printf("got client from pool for %s\n", name)
	return i, true
}

func new(ctx context.Context, opts ...options.Option) (*client, error) {
	var c client
	var err error
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: c.opts.ImpersonatedServiceAccountEmail,
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/compute"},
	})
	if err != nil {
		return nil, err
	}

	c.gce, err = gce.NewService(ctx, &ts, options.WithProject(c.opts.Project))
	if err != nil {
		return nil, err
	}
	c.cloudSql, err = cloudsql.NewService(ctx, &ts, options.WithProject(c.opts.Project))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) GCE() gce.Interface {
	return c.gce
}

func (c *client) CloudSQL() cloudsql.Interface {
	return c.cloudSql
}
