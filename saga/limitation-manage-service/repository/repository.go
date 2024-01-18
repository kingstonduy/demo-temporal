package repository

import (
	"errors"
	"fmt"
	shared "kingstonduy/demo-temporal/saga"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = shared.POSTGRES_HOST
	port     = shared.POSTGRES_PORT
	user     = shared.POSTGRES_USER
	password = shared.POSTGRES_PASSWORD
	dbname   = shared.POSTGRES_DBNAME
)

func GetConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

func DeleteUser[K any](db *gorm.DB, entity K) error {
	result := db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser[K any](db *gorm.DB, entity K) error {
	result := db.Save(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID[K any](db *gorm.DB, id string, result *K) error {
	err := db.Where("account_id = ?", id).First(&result).Error
	if err != nil {
		return errors.New("Cannot find account")
	}
	return nil
}
