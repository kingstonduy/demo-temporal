package config

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	MoneyTransfer struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"money-transfer"`
	Limit struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"limit"`
	T24 struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"t24"`
	NapasMoney struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"napas-money"`
	NapasAccount struct {
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"napas-account"`
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
}

var config *Config = nil

func GetConfig() *Config {
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

func GetDB() (db *gorm.DB, err error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		GetConfig().Database.Postgres.Host,
		GetConfig().Database.Postgres.Port,
		GetConfig().Database.Postgres.User,
		GetConfig().Database.Postgres.DBName,
		GetConfig().Database.Postgres.Password,
	)

	db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
		return nil, err
	}

	return
}

func GetAMQPChannel() *amqp.Channel {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		GetConfig().RabbitMQ.User,
		GetConfig().RabbitMQ.Password,
		GetConfig().RabbitMQ.Host,
		GetConfig().RabbitMQ.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ, %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening a channel, %s", err)
	}

	return ch
}
