package rds

import "github.com/aws/aws-sdk-go-v2/service/rds/types"

type DBInstance struct {
	types.DBInstance
}

func (i *DBInstance) hasClusterMembership() bool {
	return i.DBClusterIdentifier != nil
}

func (i *DBInstance) hasReadReplicas() bool {
	return i.ReadReplicaDBInstanceIdentifiers != nil
}

func (i *DBInstance) isReplica() bool {
	return i.ReadReplicaSourceDBClusterIdentifier != nil
}
