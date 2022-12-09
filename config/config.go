package config

import (
	"os"
	"strconv"
)

type Config struct {
	DbHost             string
	DbPort             string
	DbName             string
	DbUser             string
	DbPassword         string
	AppPort            string
	GETMethodPerMinute int
}

func InitConfig() *Config {
	config := Config{}
	threshold, _ := strconv.Atoi(os.Getenv("APP_GET_LIMIT"))
	config.DbName = os.Getenv("DB_APP_NAME")
	config.DbHost = os.Getenv("DB_APP_HOST")
	config.DbPort = os.Getenv("DB_APP_PORT")
	config.DbUser = os.Getenv("DB_APP_USER")
	config.DbPassword = os.Getenv("DB_APP_PASSWORD")
	config.AppPort = os.Getenv("APP_PORT")
	config.GETMethodPerMinute = threshold
	return &config
}
