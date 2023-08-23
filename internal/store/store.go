package store

import (
	"github.com/nduyphuong/gorya/internal/models"
)

//go:generate mockery --name Interface --output fake --outpkg fake
type Interface interface {
	SavePolicy(policy models.Policy) error
	GetPolicyByName(name string) (*models.Policy, error)
	GetPolicyBySchedule(name string) (*[]models.Policy, error)
	ListPolicy() (*[]models.Policy, error)
	DeletePolicy(name string) error
	SaveSchedule(schedule models.ScheduleModel) error
	GetSchedule(name string) (*models.ScheduleModel, error)
	ListSchedule() (*[]models.ScheduleModel, error)
	DeleteSchedule(name string) error
	SaveAKSCluster(cluster azure.AKSCluster) error
	GetAKSCluster(name string) (*azure.AKSCluster, error)
}
