package store

import (
	"fmt"
	"github.com/nduyphuong/gorya/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewCloudSQLDB(host string, user string, password string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Policy{},
		&models.ScheduleModel{},
	); err != nil {
		return db, err
	}
	return db, nil
}

