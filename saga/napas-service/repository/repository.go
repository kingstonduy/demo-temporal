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

func CreateEntity(db *gorm.DB, entity shared.NapasEntity) error {
	result := db.Create(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteUser(db *gorm.DB, entity shared.NapasEntity) error {
	result := db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser(db *gorm.DB, entity shared.NapasEntity) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID(db *gorm.DB, id string) (shared.NapasEntity, error) {
	var entity shared.NapasEntity
	result := db.First(&entity, id)
	if result.Error != nil {
		return shared.NapasEntity{}, result.Error
	}
	return entity, nil
}

func ExecuteQuery(db *gorm.DB, query string) (shared.NapasEntity, error) {
	var result shared.NapasEntity
	err := db.Raw(query).Scan(&result)
	if err != nil {
		return shared.NapasEntity{}, errors.New("Error when execute query")
	}
	return result, nil
}
