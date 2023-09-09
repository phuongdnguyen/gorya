package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type RdsClient struct {
	rds.Client
}

func NewRdsClientFromConfig(config aws.Config) *RdsClient {
	return &RdsClient{
		*rds.NewFromConfig(config),
	}
}

func (c *RdsClient) fetchAlldescribeGlobalClusters(ctx context.Context) ([]types.GlobalCluster, error) {
	globalDbClusters := []types.GlobalCluster{}
	describeGDBCOut, err := c.DescribeGlobalClusters(ctx, &rds.DescribeGlobalClustersInput{})
	if err != nil {
		return nil, err
	}
	globalDbClusters = append(globalDbClusters, describeGDBCOut.GlobalClusters...)
	for describeGDBCOut.Marker != nil {
		describeGDBCOut, err := c.DescribeGlobalClusters(ctx, &rds.DescribeGlobalClustersInput{
			Marker: describeGDBCOut.Marker,
		})
		if err != nil {
			return nil, err
		}
		globalDbClusters = append(globalDbClusters, describeGDBCOut.GlobalClusters...)
	}
	return globalDbClusters, nil
}

func (c *RdsClient) fetchAllDBCluster(ctx context.Context) ([]types.DBCluster, error) {
	dbClusters := []types.DBCluster{}
	describeDBClusterOut, err := c.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{})
	if err != nil {
		return nil, err
	}
	dbClusters = append(dbClusters, describeDBClusterOut.DBClusters...)
	for describeDBClusterOut.Marker != nil {
		describeDBClusterOut, err := c.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{Marker: describeDBClusterOut.Marker})
		if err != nil {
			return nil, err
		}
		dbClusters = append(dbClusters, describeDBClusterOut.DBClusters...)
	}
	return dbClusters, nil
}

func (c *RdsClient) fetchAllDBInstance(ctx context.Context) ([]types.DBInstance, error) {
	dbInstances := []types.DBInstance{}
	describeDBInstanceOut, err := c.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}
	dbInstances = append(dbInstances, describeDBInstanceOut.DBInstances...)
	for describeDBInstanceOut.Marker != nil {
		describeDBInstanceOut, err := c.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
			Marker: describeDBInstanceOut.Marker,
		})
		if err != nil {
			return nil, err
		}
		dbInstances = append(dbInstances, describeDBInstanceOut.DBInstances...)
	}
	return dbInstances, nil
}
