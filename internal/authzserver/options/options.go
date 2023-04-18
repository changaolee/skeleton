// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package options

import (
	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/app"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/spf13/pflag"
)

type Options struct {
	RPCServer               string                             `json:"rpcserver"      mapstructure:"rpcserver"`
	ClientCA                string                             `json:"client-ca-file" mapstructure:"client-ca-file"`
	GenericServerRunOptions *genoptions.ServerRunOptions       `json:"server"         mapstructure:"server"`
	InsecureServing         *genoptions.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	SecureServing           *genoptions.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	RedisOptions            *genoptions.RedisOptions           `json:"redis"          mapstructure:"redis"`
	Log                     *log.Options                       `json:"log"            mapstructure:"log"`
}

// NewOptions 使用默认参数创建一个 options 对象.
func NewOptions() *Options {
	o := Options{
		RPCServer:               "127.0.0.1:8081",
		ClientCA:                "",
		GenericServerRunOptions: genoptions.NewServerRunOptions(),
		InsecureServing:         genoptions.NewInsecureServingOptions(),
		SecureServing:           genoptions.NewSecureServingOptions(),
		RedisOptions:            genoptions.NewRedisOptions(),
		Log:                     log.NewOptions(),
	}
	return &o
}

func (o *Options) Flags() (fss app.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Log.AddFlags(fss.FlagSet("log"))

	o.addMiscFlags(fss.FlagSet("misc"))

	return fss
}

func (o *Options) addMiscFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
