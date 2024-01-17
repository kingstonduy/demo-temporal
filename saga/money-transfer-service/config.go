package money_transfer_service

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APPLICATION struct {
		NAME         string `yaml:"name"`
		CONTEXT_PATH string `yaml:"context-path"`
	} `yaml:"application"`

	MONEY_TRANSFER_SERVICE struct {
		NAME          string `yaml:"name"`
		WORKFLOW_NAME string `yaml:"workflow-name"`
		TASK_QUEUE    string `yaml:"task-queue"`
		PORT          string `yaml:"port"`
	} `yaml:"money_transfer_service"`

	LIMITATION_MANAGE_SERVICE struct {
		NAME          string `yaml:"name"`
		WORKFLOW_NAME string `yaml:"workflow-name"`
		TASK_QUEUE    string `yaml:"task-queue"`
		PORT          string `yaml:"port"`
	} `yaml:"limitation_manage_service:"`

	T24_SERVICE struct {
		NAME          string `yaml:"name"`
		WORKFLOW_NAME string `yaml:"workflow-name"`
		TASK_QUEUE    string `yaml:"task-queue"`
		PORT          string `yaml:"port"`
	} `yaml:"t24_service"`

	NAPAS_SERVICE struct {
		NAME          string `yaml:"name"`
		WORKFLOW_NAME string `yaml:"workflow-name"`
		TASK_QUEUE    string `yaml:"task-queue"`
		PORT          string `yaml:"port"`
	} `yaml:"napas_service"`
}

var (
	instance *Config
	once     sync.Once
)

func (c Config) GetInstance() *Config {
	once.Do(func() {
		instance = Init()
	})
	return instance
}

func Init() (config *Config) {
	configPath := GetConfigPath()
	var err error
	config, err = NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	temp := strings.Split(cwd, "/")
	temp1 := strings.Join(temp[:len(temp)-2], "/")
	path := filepath.Join(temp1, "resources/application.yml")
	return path
}
