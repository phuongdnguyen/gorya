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

//go:generate mockery --name Interface
type Interface interface {
	ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error)
}

type client struct {
	rds *RdsClient
}

func NewFromConfig(cfg aws.Config) *client {
	return &client{
		rds: NewRdsClientFromConfig(cfg),
	}
}

func (c *client) ChangeStatus(ctx context.Context, to int, tagKey string, tagValue string) (err error) {
	logger := logging.LoggerFromContext(ctx)
	if to != constants.OffStatus && to != constants.OnStatus {
		return ErrInvalidResourceStatus
	}
	seens := map[*string]types.GlobalCluster{}
	gdbClusters, err := c.rds.fetchAlldescribeGlobalClusters(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe global db clusters")
	}
	/**
	  https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-cluster-stop-start.html
	  **/
	for _, gCluster := range gdbClusters {
		for _, member := range gCluster.GlobalClusterMembers {
			seens[member.DBClusterArn] = gCluster
		}
		logger.Infof("found %s is global db cluster, skip", *gCluster.GlobalClusterIdentifier)
	}
	dbClusters, err := c.rds.fetchAllDBCluster(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe db clusters")
	}
	dbInstances, err := c.rds.fetchAllDBInstance(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "describe db instances")
	}
	for _, cluster := range dbClusters {
		if gdbIdentifier, exist := seens[cluster.DBClusterArn]; exist {
			logger.Infof("found %s is member of aurora global db %s, skip", *cluster.DBClusterArn, *gdbIdentifier.GlobalClusterIdentifier)
			continue
		}
		for _, tag := range cluster.TagList {
			if *tag.Key != tagKey {
				continue
			}
			if *tag.Value != tagValue {
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
	}
	for _, instance := range dbInstances {
		dbInstance := DBInstance{
			instance,
		}
		if dbInstance.hasClusterMembership() {
			logger.Infof("instance %s is member of cluster %s, skip", *dbInstance.DBInstanceIdentifier, *dbInstance.DBClusterIdentifier)
			continue
		}
		if dbInstance.hasReadReplicas() {
			logger.Infof("instance %s has read replicas, skip", *dbInstance.DBInstanceIdentifier)
			continue
		}
		if dbInstance.isReplica() {
			logger.Infof("instance %s is a read replica, skip", *dbInstance.DBInstanceIdentifier)
			continue
		}
		for _, tag := range dbInstance.TagList {
			if *tag.Key != tagKey {
				continue
			}
			if *tag.Value != tagValue {
				continue
			}
			/**
			  https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/custom-managing-sqlserver.html#custom-managing-sqlserver.startstop
			  **/
			if dbInstance.AutomationMode == types.AutomationModeFull {
				_, err = c.rds.ModifyDBInstance(ctx, &rds.ModifyDBInstanceInput{
					DBInstanceIdentifier: dbInstance.DBInstanceIdentifier,
					AutomationMode:       types.AutomationModeAllPaused,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("stop rds custom automation%s", *dbInstance.DBInstanceIdentifier)))
					break
				}

			}
			switch to {
			case constants.OnStatus:
				_, err = c.rds.StartDBInstance(ctx, &rds.StartDBInstanceInput{
					DBInstanceIdentifier: dbInstance.DBInstanceIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("start db instance %s", *dbInstance.DBInstanceIdentifier)))
				}
			case constants.OffStatus:
				_, err = c.rds.StopDBInstance(ctx, &rds.StopDBInstanceInput{
					DBInstanceIdentifier: dbInstance.DBInstanceIdentifier,
				})
				if err != nil {
					logger.Error(pkgerrors.Wrap(err, fmt.Sprintf("stop db instance %s", *dbInstance.DBInstanceIdentifier)))
				}
			}
		}
	}

	return nil
}
