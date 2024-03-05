package bootstrap

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (c *Config) GetDB() (db *gorm.DB, err error) {
	var dbConn *gorm.DB = nil
	if dbConn != nil {
		return dbConn, nil
	}

	fmt.Println("💡💡💡 Create connection")
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		c.Database.Postgres.Host,
		c.Database.Postgres.Port,
		c.Database.Postgres.User,
		c.Database.Postgres.DBName,
		c.Database.Postgres.Password,
	)

	dbConn, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
		return nil, err
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatalf("Error getting underlying sql.DB, %s", err)
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(50)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return dbConn, nil
}
