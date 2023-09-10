package store

import (
	"sync"

	"gorm.io/gorm/clause"

	"github.com/nduyphuong/gorya/internal/constants"
	"github.com/nduyphuong/gorya/internal/models"
	"github.com/nduyphuong/gorya/internal/os"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

var modelStorage *Storage
var muModelStorage sync.Mutex

func GetOnce() (Interface, error) {
	muModelStorage.Lock()
	defer func() {
		muModelStorage.Unlock()
	}()
	if modelStorage != nil {
		return modelStorage, nil
	}
	dbType := os.GetEnv(constants.ENV_GORYA_DB_TYPE, "sqlite")
	host := os.GetEnv(constants.ENV_GORYA_DB_HOST, "localhost:3306")
	user := os.GetEnv(constants.ENV_GORYA_DB_USER, "root")
	password := os.GetEnv(constants.ENV_GORYA_DB_PASSWORD, "root")
	dbName := os.GetEnv(constants.ENV_GORYA_DB_NAME, "gorya")
	switch dbType {
	case "mysql":
		db, err := NewMySQLDB(host, user, password, dbName)
		if err != nil {
			return nil, err
		}
		store := New(db)
		return store, nil
	default:
		db, err := NewSqliteDB()
		if err != nil {
			return nil, err
		}
		store := New(db)
		return store, nil
	}
}

func New(db *gorm.DB) Interface {
	return &Storage{
		db,
	}
}

func (c *Storage) SavePolicy(m models.Policy) error {

	if err := c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "display_name", "projects", "tags", "schedule_name", "provider"}),
	}).Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (c *Storage) GetPolicyByName(name string) (*models.Policy, error) {
	r := models.Policy{}
	query := c.db.Model(&models.Policy{}).Where("name=?", name)
	if err := query.First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) GetPolicyBySchedule(name string) (*[]models.Policy, error) {
	var r []models.Policy
	query := c.db.Model(&models.Policy{}).Where("schedule_name=?", name)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) ListPolicy() (*[]models.Policy, error) {
	var r []models.Policy
	if err := c.db.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) ListPolicyByProvider(name string) (*[]models.Policy, error) {
	var r []models.Policy
	if err := c.db.Where("name <> ?", name).Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) DeletePolicy(name string) error {
	r := models.Policy{}
	query := c.db.Model(&models.Policy{}).Where("name=?", name)
	if err := query.Delete(&r).Error; err != nil {
		return err
	}
	return nil
}

func (c *Storage) SaveSchedule(m models.ScheduleModel) error {
	if err := c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "display_name", "time_zone", "schedule"}),
	}).Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (c *Storage) GetSchedule(name string) (*models.ScheduleModel, error) {
	r := models.ScheduleModel{}
	query := c.db.Model(&models.ScheduleModel{}).Where("name = ?", name)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) ListSchedule() (*[]models.ScheduleModel, error) {
	var r []models.ScheduleModel
	if err := c.db.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Storage) DeleteSchedule(name string) error {
	r := models.ScheduleModel{}
	query := c.db.Model(&models.ScheduleModel{}).Where("name=?", name)
	if err := query.Delete(&r).Error; err != nil {
		return err
	}
	return nil
}
