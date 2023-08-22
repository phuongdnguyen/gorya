package config

import (
	"os"
)

type EKSConfig struct {
	ClusterName    string
	NodeGroupName  string
	// Add other necessary fields here
}

func LoadEKSConfigFromEnv() *EKSConfig {
	return &EKSConfig{
		ClusterName:    os.Getenv("EKS_CLUSTER_NAME"),
		NodeGroupName:  os.Getenv("EKS_NODE_GROUP_NAME"),
		// Load other necessary fields here
	}
}

