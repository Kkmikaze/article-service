package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDNS string `mapstructure:"POSTGRES_DNS"`
}

var Configs *Config

func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Configs)
	if err != nil {
		panic(err)
	}
}
