package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Target     string `mapstructure:"target"`
	Duration   int    `mapstructure:"duration"`
	File       string `mapstructure:"file"`
	ServerAddr string `mapstructure:"serverAddr"`
}

var vp *viper.Viper
var config Config

func LoadConfig(cfgFile *string) (*Config, error) {
	if *cfgFile == "" {
		return nil, nil
	}

	vp = viper.New()

	vp.SetConfigFile(*cfgFile)

	vp.SetDefault("target", "ws://localhost:8000/")
	vp.SetDefault("duration", 0)
	vp.SetDefault("file", nil)
	vp.SetDefault("ServerAddr", "localhost:8000")

	err := vp.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

func NewRecordConfig(target *string, duration int, file *string) (*Config, error) {
	if *target == "" {
		return nil, errors.New("missing target parameter")
	}
	if *file == "" {
		return nil, errors.New("missing file parameter")
	}
	config = Config{
		Target:   *target,
		Duration: duration,
		File:     *file,
	}
	return &config, nil
}

func NewPlaybackConfig(file *string, serverAddr *string) (*Config, error) {
	if *file == "" {
		return nil, errors.New("missing file parameter")
	}
	if *serverAddr == "" {
		return nil, errors.New("missing server parameter")
	}
	config = Config{
		File:       *file,
		ServerAddr: *serverAddr,
	}
	return &config, nil
}
