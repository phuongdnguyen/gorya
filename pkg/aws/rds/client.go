package rds

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	pkgerrors "github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/logging"
)

var ErrInvalidResourceStatus = errors.New("invalid resource status")

type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

type client struct {
	rds *rds.Client
}

func NewFromConfig(cfg aws.Config) (*client, error) {
	c := &client{}
	c.rds = rds.NewFromConfig(cfg)
	return c, nil
}

func (c *client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error) {
	logger := logging.LoggerFromContext(ctx)
	if to != constants.OffStatus && to != constants.OnStatus {
		return ErrInvalidResourceStatus
	}
	dbClusters, err := c.describeDBCluster(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe db clusters")
	}
	// fmt.Printf("dbClusters: %v\n", dbClusters)
	dbInstances, err := c.describeDBInstance(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe db instances")
	}
	// fmt.Printf("dbInstances: %v\n", dbInstances)
	instanceToClusterMap := map[*string]*string{}
	for _, cluster := range dbClusters {
		// fmt.Printf("cluster.DBClusterIdentifier: %v\n", *cluster.DBClusterIdentifier)
		for _, tag := range cluster.TagList {
			// fmt.Printf("tag.Key: %v\n", *tag.Key)
			// fmt.Printf("tag.Value: %v\n", *tag.Value)
			if *tag.Key != tagKey || *tag.Value != tagValue {
				continue
			}
			switch to {
			case constants.OnStatus:
				_, err = c.rds.StartDBCluster(ctx, &rds.StartDBClusterInput{
					DBClusterIdentifier: cluster.DBClusterIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("start db cluster %s", *cluster.DBClusterIdentifier)))
				}
			case constants.OffStatus:
				_, err = c.rds.StopDBCluster(ctx, &rds.StopDBClusterInput{
					DBClusterIdentifier: cluster.DBClusterIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("stop db cluster %s", *cluster.DBClusterIdentifier)))
				}
			}

		}
		for _, member := range cluster.DBClusterMembers {
			instanceToClusterMap[member.DBInstanceIdentifier] = cluster.DBClusterIdentifier
		}
	}
	for _, instance := range dbInstances {
		for _, tag := range instance.TagList {
			if *tag.Key != tagKey || *tag.Value != tagValue || instanceToClusterMap[instance.DBInstanceIdentifier] != nil {
				continue
			}
			if instance.AutomationMode == types.AutomationModeFull {
				_, err = c.rds.ModifyDBInstance(ctx, &rds.ModifyDBInstanceInput{
					DBInstanceIdentifier: instance.DBInstanceIdentifier,
					AutomationMode:       types.AutomationModeAllPaused,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("stop rds custom automation%s", *instance.DBInstanceIdentifier)))
					break
				}

			}
			switch to {
			case constants.OnStatus:
				_, err = c.rds.StartDBInstance(ctx, &rds.StartDBInstanceInput{
					DBInstanceIdentifier: instance.DBInstanceIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("start db instance %s", *instance.DBInstanceIdentifier)))
				}
			case constants.OffStatus:
				_, err = c.rds.StopDBInstance(ctx, &rds.StopDBInstanceInput{
					DBInstanceIdentifier: instance.DBInstanceIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("stop db instance %s", *instance.DBInstanceIdentifier)))
				}
			}
		}
	}

	return nil
}

func (c *client) describeDBCluster(ctx context.Context) ([]types.DBCluster, error) {
	dbClusters := []types.DBCluster{}
	describeDBClusterOut, err := c.rds.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{})
	if err != nil {
		return nil, err
	}
	dbClusters = append(dbClusters, describeDBClusterOut.DBClusters...)
	for describeDBClusterOut.Marker != nil {
		describeDBClusterOut, err := c.rds.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{Marker: describeDBClusterOut.Marker})
		if err != nil {
			return nil, err
		}
		dbClusters = append(dbClusters, describeDBClusterOut.DBClusters...)
	}
	return dbClusters, nil
}

func (c *client) describeDBInstance(ctx context.Context) ([]types.DBInstance, error) {
	dbInstances := []types.DBInstance{}
	describeDBInstanceOut, err := c.rds.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}
	dbInstances = append(dbInstances, describeDBInstanceOut.DBInstances...)
	for describeDBInstanceOut.Marker != nil {
		describeDBInstanceOut, err := c.rds.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
			Marker: describeDBInstanceOut.Marker,
		})
		if err != nil {
			return nil, err
		}
		dbInstances = append(dbInstances, describeDBInstanceOut.DBInstances...)
	}
	return dbInstances, nil
}
