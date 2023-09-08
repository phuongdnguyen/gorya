package cloudsql

import (
	"context"
	"errors"

	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"github.com/nduyphuong/gorya/pkg/gcp/utils"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	sql "google.golang.org/api/sqladmin/v1beta4"
)

var ErrInvalidResourceStatus = errors.New("invalid resource status")

type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

type client struct {
	sql  *sql.Service
	opts options.Options
}

func NewService(ctx context.Context, ts *oauth2.TokenSource, opts ...options.Option) (*client, error) {
	var err error
	c := &client{}
	for _, o := range opts {
		o.Apply(&c.opts)
	}
	c.sql, err = sql.NewService(ctx, option.WithTokenSource(*ts))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) error {
	if to != constants.OffStatus && to != constants.OnStatus {
		return ErrInvalidResourceStatus
	}
	var action string
	if to == constants.OffStatus {
		action = "NEVER"
	}
	if to == constants.OnStatus {
		action = "ALWAYS"
	}
	tagFilter := utils.GetCloudSqlFilter(tagKey, tagValue)
	instancesListResp, err := c.sql.Instances.List(c.opts.Project).Filter(tagFilter).Do()
	if err != nil {
		return pkgerrors.Wrap(err, "list instances")
	}
	eg, _ := errgroup.WithContext(ctx)
	for _, instance := range instancesListResp.Items {
		instance := instance
		eg.Go(func() error {
			rb := &sql.DatabaseInstance{
				Settings: &sql.Settings{
					ActivationPolicy: action,
				},
			}
			_, err := c.sql.Instances.Patch(c.opts.Project, instance.Name, rb).Do()
			if err != nil && !googleapi.IsNotModified(err) {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}
