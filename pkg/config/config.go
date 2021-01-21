package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config interface {
	ParamConfig() ParamConfig
}

type appConfig struct {
	paramConfig ParamConfig
}

func (ac appConfig) ParamConfig() ParamConfig {
	return ac.paramConfig
}

func NewConfig(configFile string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return appConfig{
		paramConfig: newParamConfig(),
	}
}
