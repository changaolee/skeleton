package apiserver

import (
	"github.com/changaolee/skeleton/internal/apiserver/options"
	"github.com/changaolee/skeleton/pkg/app"
)

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("SKT API Server",
		basename,
		app.WithOptions(opts),
	)
	return application
}
