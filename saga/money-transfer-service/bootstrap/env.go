package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerHost      string `mapstructure:"SERVER_HOST"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPass          string `mapstructure:"DB_PASS"`
	DBName          string `mapstructure:"DB_NAME"`
	ContextTimeout  int    `mapstructure:"CONTEXT_TIMEOUT"`
	NapasUrl        string `mapstructure:"NAPAS_URL"`
	T24Url          string `mapstructure:"T24_URL"`
	LimitServiceUrl string `mapstructure:"LIMIT_SERVICE_URL"`
	TaskQueue       string `mapstructure:"TASK_QUEUE"`
	Workflow        string `mapstructure:"WORKFLOW"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
