package config

import "github.com/changaolee/skeleton/internal/apiserver/options"

type Config struct {
	*options.Options
}

// CreateConfigFromOptions 基于给定的选项创建应用配置.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
