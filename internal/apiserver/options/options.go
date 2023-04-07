// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package options

import (
	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/app"
	"github.com/changaolee/skeleton/pkg/log"
)

type Options struct {
	GenericServerRunOptions *genoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing         *genoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQLOptions            *genoptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	Log                     *log.Options                       `json:"log"      mapstructure:"log"`
}

// NewOptions 使用默认参数创建一个 options 对象.
func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: genoptions.NewServerRunOptions(),
		InsecureServing:         genoptions.NewInsecureServingOptions(),
		SecureServing:           genoptions.NewSecureServingOptions(),
		MySQLOptions:            genoptions.NewMySQLOptions(),
		Log:                     log.NewOptions(),
	}
	return &o
}

func (o *Options) Flags() (fss app.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.Log.AddFlags(fss.FlagSet("log"))

	return fss
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
