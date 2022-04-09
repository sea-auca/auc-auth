package config

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
)

//Main application configuration
type AppConfig struct {
	// Http server configuration block. Has basic host and port params with default timeout values
	Server struct {
		Port        int
		ReadTimeout int
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
	// configuration for the email sender
	Email struct {
		Host string
		Port int
	}
	Service struct {
		VerificationPrefix string
	}
}

var conf AppConfig

func Init(lg *zap.Logger) {
	defer func() {
		if r := recover(); r != nil {
			lg.Fatal("Failed to initialise configuration", zap.Any("error", r))
		}
	}()
	conf.Server.Port = readIntField("SERVER_PORT")
	conf.Server.ReadTimeout = readIntField("SERVER_READTIMEOUT")
	conf.Email.Host = readField("EMAIL_HOST")
	conf.Email.Port = readIntField("EMAIL_PORT")
	conf.Database.Host = readField("DATABASE_HOST")
	conf.Database.Port = readField("DATABASE_PORT")
	conf.Database.User = readField("DATABASE_USER")
	conf.Database.Password = readField("DATABASE_PASSWORD")
	conf.Database.Database = readField("DATABASE_NAME")
	conf.Service.VerificationPrefix = readField("SERVICE_VERIFY_URL")
}

func Config() AppConfig {
	return conf
}

func readField(f string) string {
	val := os.Getenv(f)
	if val == "" {
		panic(fmt.Errorf("key %s is not present. panic", f))
	}
	return val
}

func readIntField(f string) int {
	val := os.Getenv(f)
	if val == "" {
		panic(fmt.Errorf("key %s is not present. panic", f))
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("key value can not be parsed, key: %s", f))
	}
	return res
}
