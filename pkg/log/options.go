// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package log

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagOutputPaths       = "log.output-paths"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options 包含与日志相关的配置项.
type Options struct {
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`     // 是否禁止 caller，如果开启会在日志中显示调用日志所在的文件和行号
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // 是否禁止在 panic 及以上级别打印堆栈信息
	Level             string   `json:"level"              mapstructure:"level"`              // 指定日志级别，可选值：debug, info, warn, error, dpanic, panic, fatal
	Format            string   `json:"format"             mapstructure:"format"`             // 指定日志格式，可选值：console, json
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`       // 指定日志输出位置
}

// NewOptions 创建一个带有默认参数的 Options 对象.
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		OutputPaths:       []string{"stdout"},
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
}

func (o *Options) Validate() []error {
	var errs []error

	// 检查日志级别是否合法
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	// 检查日志格式是否合法
	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}
