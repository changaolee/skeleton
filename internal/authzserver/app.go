// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package authzserver

import (
	"github.com/changaolee/skeleton/internal/authzserver/config"
	"github.com/changaolee/skeleton/internal/authzserver/options"
	"github.com/changaolee/skeleton/pkg/app"
	"github.com/changaolee/skeleton/pkg/log"
)

const commandDesc = `The SKT Authorization Server to run ladon policies which can protecting your resources.`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("SKT Authorization Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Sync()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
