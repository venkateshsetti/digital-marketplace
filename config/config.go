package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MongoConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type CoingeckoConfig struct {
	Url string `json:"url"`
}

type AppConfig struct {
	Server    ServerConfig
	Redis     RedisConfig
	Mongo     MongoConfig
	Coingecko CoingeckoConfig
}

func LoadConfig() *AppConfig {
	var configuration *AppConfig
	var configName string
	configName = "default_config" // single config file
	viper.SetConfigName(configName)
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err.Error())
	}
	err := viper.MergeInConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = viper.UnmarshalExact(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return configuration
}
