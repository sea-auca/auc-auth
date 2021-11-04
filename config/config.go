package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	IsDevelopmentConfig bool `yaml:"is_dev"`
	ServerConfig        struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Read  int
			Write int
		}
	} `yaml:"server"`
	DatabaseConfig struct {
		Host            string
		Port            string
		User            string
		Password        string
		ConnectionLimit int
	} `yaml:"database"`
}

var ErrNoConfigFile = errors.New("config file does not exist")
var ErrFileIssue = errors.New("config file can not be opened")
var ErrParseIssue = errors.New("configuration file can not be parsed")

func ReadConfig(configPath *string) (*AppConfig, error) {
	path := "./config/config.dev.yaml"
	if val, exists := os.LookupEnv("PROD_CONFIG"); exists {
		path = val
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrNoConfigFile
	}

	config := &AppConfig{}

	file, err := os.Open(path)
	if err != nil {
		return nil, ErrFileIssue
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("%s. Error: %s", ErrParseIssue.Error(), err.Error())
	}

	return config, nil
}
