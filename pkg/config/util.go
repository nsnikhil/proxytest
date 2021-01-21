package config

import (
	"github.com/spf13/viper"
)

func getString(config string, defaultVal ...string) string {
	if len(defaultVal) > 0 {
		viper.SetDefault(config, defaultVal[0])
	}

	return viper.GetString(config)
}
