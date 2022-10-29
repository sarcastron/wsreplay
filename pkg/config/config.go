package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Target   string `mapstructure:"target"`
	Duration int    `mapstructure:"duration"`
}

var defaultConfig Config = Config{
	"ws://localhost:8000/",
	0,
}

var vp *viper.Viper
var config Config

func loadConfig(cfgFile string) (Config, error) {
	vp = viper.New()

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

func GetConfig(cfgFile string) (Config, error) {
	var err error = nil
	if config.Target == "" {
		config, err = loadConfig(cfgFile)
	}
	return config, err
}