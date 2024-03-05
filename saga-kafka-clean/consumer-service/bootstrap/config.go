package bootstrap

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	MoneyTransfer struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
		Kafka struct {
			Topic struct {
				In  string `mapstructure:"in"`
				Out string `mapstructure:"out"`
			} `mapstructure:"topic"`
		} `mapstructure:"kafka"`
	} `mapstructure:"money-transfer"`
	Limit struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
		Kafka struct {
			Topic struct {
				In  string `mapstructure:"in"`
				Out string `mapstructure:"out"`
			} `mapstructure:"topic"`
		} `mapstructure:"kafka"`
	} `mapstructure:"limit"`
	T24 struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
		Kafka struct {
			Topic struct {
				In  string `mapstructure:"in"`
				Out string `mapstructure:"out"`
			} `mapstructure:"topic"`
		} `mapstructure:"kafka"`
	} `mapstructure:"t24"`
	Napas struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
		Kafka struct {
			Topic struct {
				In  string `mapstructure:"in"`
				Out string `mapstructure:"out"`
			} `mapstructure:"topic"`
		} `mapstructure:"kafka"`
	} `mapstructure:"napas"`
	Database struct {
		Postgres struct {
			Host     string `mapstructure:"host"`
			Port     string `mapstructure:"port"`
			DBName   string `mapstructure:"dbname"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
		} `mapstructure:"postgres"`
	} `mapstructure:"database"`
	Temporal struct {
		Host      string `mapstructure:"host"`
		Port      string `mapstructure:"port"`
		TaskQueue string `mapstructure:"taskqueue"`
		Workflow  string `mapstructure:"workflow"`
	} `mapstructure:"temporal"`
	RabbitMQ struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"rabbitmq"`
	Kafka struct {
		BootstrapServer struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"bootstrap-server"`
	} `mapstructure:"kafka"`
}

var config *Config = nil

func NewConfig() *Config {
	var cfg *viper.Viper
	if cfg == nil {
		cfg = viper.New()
		cfg.SetConfigType("yml")
		cfg.SetConfigFile("../application.yml")

		err := cfg.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	err := cfg.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode stringo struct, %v", err)
	}

	return config
}

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

func (c *Config) GetAMQPConnection() *amqp.Connection {
	var conn *amqp.Connection = nil

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ, %s", err)
	}

	return conn
}

func (c *Config) GetTemporalClient() *client.Client {
	var temporalClient client.Client = nil
	temporalClient, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", c.Temporal.Host, c.Temporal.Port),
	})
	if err != nil {
		log.Fatalf("unable to create client, %v", err)
	}

	return &temporalClient
}
