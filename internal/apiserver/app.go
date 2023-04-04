package apiserver

import (
	"github.com/changaolee/skeleton/internal/apiserver/config"
	"github.com/changaolee/skeleton/internal/apiserver/options"
	"github.com/changaolee/skeleton/pkg/app"
	"github.com/changaolee/skeleton/pkg/log"
)

const commandDesc = `The SKT API Server validates and configures data
for the api objects which include users, policies, secrets, and others.

The API Server services REST operations to do the api objects management.`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("SKT API Server",
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
