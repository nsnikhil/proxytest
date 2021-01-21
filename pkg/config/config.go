package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config interface {
	ParamConfig() ParamConfig
	RateLimitConfig() RateLimitConfig
}

type appConfig struct {
	paramConfig     ParamConfig
	rateLimitConfig RateLimitConfig
}

func (ac appConfig) ParamConfig() ParamConfig {
	return ac.paramConfig
}

func (ac appConfig) RateLimitConfig() RateLimitConfig {
	return ac.rateLimitConfig
}

func NewConfig(configFile string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return appConfig{
		paramConfig:     newParamConfig(),
		rateLimitConfig: newRateLimitConfig(),
	}
}
