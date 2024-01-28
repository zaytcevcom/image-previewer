package main

import "github.com/spf13/viper"

type Config struct {
	Logger LoggerConf
	Cache  CacheConf
}

type LoggerConf struct {
	Level string
}

type CacheConf struct {
	Capacity int
}

func LoadConfig(path string) (Config, error) {
	config := Config{}

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
