package rds

import (
	"context"
	"errors"

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
		return pkgerrors.Wrap(err, "describe dbclusters")
	}
	dbInstances, err := c.describeDBInstance(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe dbinstances")
	}
	instanceToClusterMap := map[*string]*string{}
	for _, cluster := range dbClusters {
		for _, tag := range cluster.TagList {
			if *tag.Key != tagKey || *tag.Value == tagValue {
				continue
			}
			_, err = c.rds.StopDBCluster(ctx, &rds.StopDBClusterInput{
				DBClusterIdentifier: cluster.DBClusterIdentifier,
			})
			if err != nil {
				logger.Errorf("stop db cluster %v", cluster.DBClusterIdentifier)
			}
		}
		for _, member := range cluster.DBClusterMembers {
			instanceToClusterMap[member.DBInstanceIdentifier] = cluster.DBClusterIdentifier
		}
	}
	for _, instance := range dbInstances {
		for _, tag := range instance.TagList {
			if *tag.Key != tagKey || *tag.Value == tagValue || instanceToClusterMap[instance.DBInstanceIdentifier] != nil {
				continue
			}
			_, err = c.rds.StopDBInstance(ctx, &rds.StopDBInstanceInput{
				DBInstanceIdentifier: instance.DBInstanceIdentifier,
			})
			if err != nil {
				logger.Errorf("stop db instance %v", instance.DBInstanceIdentifier)
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
