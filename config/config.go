package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(configPath())
	viper.SetConfigName(configName())
	viper.SetConfigType(configType())
}

func Parse() error {
	return viper.ReadInConfig()
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func configPath() string {
	if configPath := os.Getenv("CONFIG_PATH"); configPath == "" {
		return "."
	} else {
		return configPath
	}
}

func configName() string {
	if configName := os.Getenv("CONFIG_NAME"); configName == "" {
		return "config"
	} else {
		return configName
	}
}

func configType() string {
	if configType := os.Getenv("CONFIG_TYPE"); configType == "" {
		return "yaml"
	} else {
		return configType
	}
}

func Setup() {
	log.Println("Setting up configuration")
	if err := Parse(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
		return
	}
}
