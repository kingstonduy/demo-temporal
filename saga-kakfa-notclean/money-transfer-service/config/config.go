package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MoneyTransfer struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"money-transfer"`
	Limit struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"limit"`
	T24 struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"t24"`
	Napas struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"napas"`
	Database struct {
		Postgres struct {
			Host     string `mapstructure:"host"`
			Port     string `mapstructure:"port"`
			DBName   string `mapstructure:"dbname"`
			Password string `mapstructure:"password"`
		} `mapstructure:"postgres"`
	} `mapstructure:"database"`
	Temporal struct {
		Host      string `mapstructure:"host"`
		Port      string `mapstructure:"port"`
		TaskQueue string `mapstructure:"taskqueue"`
		Workflow  string `mapstructure:"workflow"`
	} `mapstructure:"temporal"`
	Kafka struct {
		Index           string `mapstructure:"index"`
		BootstrapServer struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"bootstrap-server"`
	} `mapstructure:"kafka"`
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
