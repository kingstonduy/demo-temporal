package bootstrap

import (
	"log"

	"github.com/spf13/viper"
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

func GetConfig() *Config {
	var config *Config
	var cfg *viper.Viper

	cfg = viper.New()
	cfg.SetConfigType("yml")
	cfg.SetConfigFile("./application.yml")

	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = cfg.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode string struct, %v", err)
	}

	return config
}
