package shared

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = POSTGRES_HOST
	port     = POSTGRES_PORT
	user     = POSTGRES_USER
	password = POSTGRES_PASSWORD
	dbname   = POSTGRES_DBNAME
)

func GetConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}

func CreateEntity[K any](db *gorm.DB, entity K) error {
	result := db.Create(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteEntity[K any](db *gorm.DB, entity K) error {
	result := db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateEntity[K any](db *gorm.DB, entity K) error {
	result := db.Save(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetEntityByID[K any](db *gorm.DB, id string, result *K) error {
	err := db.Where("account_id = ?", id).First(&result).Error
	if err != nil {
		return errors.New("Cannot find account")
	}
	return nil
}
