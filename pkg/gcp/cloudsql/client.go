package cloudsql

import (
	"context"
	"errors"
	"sync"

	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/logging"
	"github.com/nduyphuong/gorya/pkg/gcp/options"
	"github.com/nduyphuong/gorya/pkg/gcp/utils"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/oauth2"
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
	logger := logging.LoggerFromContext(ctx)
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
	replicasToInstance := map[string]string{}
	var wg sync.WaitGroup
	wg.Add(len(instancesListResp.Items))
	for _, instance := range instancesListResp.Items {
		instance := instance
		if len(instance.ReplicaNames) > 0 {
			for _, replName := range instance.ReplicaNames {
				replicasToInstance[replName] = instance.Name
			}
		}
		if _, exist := replicasToInstance[instance.Name]; exist {
			// instance is replica
			continue
		}
		go func() {
			defer wg.Done()
			rb := &sql.DatabaseInstance{
				Settings: &sql.Settings{
					ActivationPolicy: action,
				},
			}
			_, err := c.sql.Instances.Patch(c.opts.Project, instance.Name, rb).Do()
			if err != nil && !googleapi.IsNotModified(err) {
				logger.Errorf("patch instance %s", instance.Name)
			}
		}()
	}
	return nil
}
