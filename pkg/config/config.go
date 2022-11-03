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

var ErrMissingFileParam error = errors.New("missing file parameter")
var ErrMissingTargetParam error = errors.New("missing target parameter")
var ErrMissingServerParam error = errors.New("missing server parameter")

func LoadConfig(cfgFile *string) (*Config, error) {
	if *cfgFile == "" {
		return nil, nil
	}

	vp = viper.New()

	vp.SetConfigFile(*cfgFile)

	vp.SetDefault("duration", 0)
	vp.SetDefault("ServerAddr", ":8000")

	err := vp.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return &Config{}, err
	}

	if config.File == "" {
		return nil, ErrMissingFileParam
	}

	if config.Target == "" {
		return nil, ErrMissingTargetParam
	}

	return &config, nil
}

func NewRecordConfig(target *string, duration int, file *string) (*Config, error) {
	if *target == "" {
		return nil, ErrMissingTargetParam
	}
	if *file == "" {
		return nil, ErrMissingFileParam
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
		return nil, ErrMissingFileParam
	}
	if *serverAddr == "" {
		return nil, ErrMissingServerParam
	}
	config = Config{
		File:       *file,
		ServerAddr: *serverAddr,
	}
	return &config, nil
}
