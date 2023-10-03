package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig(file string) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while loading config: %v", err)
	}
}
