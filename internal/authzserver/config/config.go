// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package config

import "github.com/changaolee/skeleton/internal/authzserver/options"

type Config struct {
	*options.Options
}

// CreateConfigFromOptions 基于给定的选项创建应用配置.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
