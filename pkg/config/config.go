package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config interface {
	Env() string
	ParamConfig() ParamConfig
	RateLimitConfig() RateLimitConfig
	HTTPClientConfig() HTTPClientConfig
	HTTPServerConfig() HTTPServerConfig
	LogConfig() LogConfig
	LogFileConfig() LogFileConfig
}

type appConfig struct {
	env              string
	paramConfig      ParamConfig
	rateLimitConfig  RateLimitConfig
	httpClientConfig HTTPClientConfig
	httpServerConfig HTTPServerConfig
	logConfig        LogConfig
	logFileConfig    LogFileConfig
}

func (ac appConfig) Env() string {
	return ac.env
}

func (ac appConfig) ParamConfig() ParamConfig {
	return ac.paramConfig
}

func (ac appConfig) RateLimitConfig() RateLimitConfig {
	return ac.rateLimitConfig
}

func (ac appConfig) HTTPClientConfig() HTTPClientConfig {
	return ac.httpClientConfig
}

func (ac appConfig) HTTPServerConfig() HTTPServerConfig {
	return ac.httpServerConfig
}

func (ac appConfig) LogConfig() LogConfig {
	return ac.logConfig
}

func (ac appConfig) LogFileConfig() LogFileConfig {
	return ac.logFileConfig
}

func NewConfig(configFile string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return appConfig{
		env:              getString("ENV"),
		paramConfig:      newParamConfig(),
		rateLimitConfig:  newRateLimitConfig(),
		httpClientConfig: newHTTPClientConfig(),
		httpServerConfig: newHTTPServerConfig(),
		logConfig:        newLogConfig(),
		logFileConfig:    newLogFileConfig(),
	}
}
