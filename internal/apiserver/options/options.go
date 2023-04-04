package options

import (
	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/app"
	"github.com/changaolee/skeleton/pkg/log"
)

type Options struct {
	MySQLOptions *genoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	Log          *log.Options             `json:"log"   mapstructure:"log"`
}

func (o *Options) Flags() (fss app.NamedFlagSets) {
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

// NewOptions 使用默认参数创建一个 options 对象.
func NewOptions() *Options {
	o := Options{
		MySQLOptions: genoptions.NewMySQLOptions(),
		Log:          log.NewOptions(),
	}
	return &o
}
