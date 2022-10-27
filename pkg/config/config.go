package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Target   string `mapstructure:"target"`
	Duration int    `mapstructure:"duration"`
}

var defaultConfig Config = Config{
	"localhost:8000",
	0,
}

var vp *viper.Viper

func LoadConfig(cfgFile string) (Config, error) {
	vp = viper.New()
	var config Config

	if cfgFile != "" {
		// Use config file from the flag.
		vp.SetConfigFile(cfgFile)
	} else {
		vp.AddConfigPath("")
		vp.SetConfigName("config")
		vp.SetConfigType("yaml")
	}

	vp.SetDefault("target", defaultConfig.Target)
	vp.SetDefault("duration", defaultConfig.Duration)

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
