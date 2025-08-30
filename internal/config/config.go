package config

import "github.com/spf13/viper"

type Config struct {
	Port     int          `mapstructure:"port"`
	Database DatabaseConf `mapstructure:"database"`
	Token    TokenConf    `mapstructure:"token"`
}
type DatabaseConf struct {
	Url string
}
type TokenConf struct {
	Secret string // as we use md5 for jwt token a string secret is enought
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetDefault("port", "8080")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
