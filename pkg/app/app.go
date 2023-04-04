package app

import (
	"fmt"
	"os"

	"github.com/changaolee/skeleton/pkg/log"
	"github.com/changaolee/skeleton/pkg/version"
	"github.com/changaolee/skeleton/pkg/version/verflag"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var progressMessage = color.GreenString("==>")

// App 是命令行应用程序的主体结构，建议使用 app.NewApp() 函数创建应用.
type App struct {
	basename    string               // 二进制文件名
	name        string               // 应用名
	description string               // 应用描述
	options     CliOptions           // 命令行选项参数
	runFunc     RunFunc              // 应用程序启动的回调函数
	silence     bool                 // 是否打印执行信息
	noVersion   bool                 // 是否不显示版本信息
	noConfig    bool                 // 是否未指定配置文件
	commands    []*Command           // 应用程序包含的子命令
	args        cobra.PositionalArgs // 预期参数
	cmd         *cobra.Command       // 根命令，用于启动应用
}

// Option 定义用于初始化 App 结构的可选参数.
type Option func(*App)

func WithOptions(opts CliOptions) Option {
	return func(a *App) {
		a.options = opts
	}
}

// RunFunc 定义应用程序启动的回调函数.
type RunFunc func(basename string) error

// NewApp 基于给定的应用名、二进制文件名以及一些其他选项，创建一个应用程序实例.
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}
	for _, opt := range opts {
		opt(a)
	}
	a.buildCommand()

	return a
}

// Run 用于启动应用程序.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           a.basename,
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.basename))
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}
	addHelpFlag(cmd.Name(), namedFlagSets.FlagSet("global"))
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	addCmdTemplate(&cmd, namedFlagSets)
	a.cmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	PrintFlags(cmd.Flags())

	if !a.noVersion {
		// 打印应用程序版本信息
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// 启动应用程序
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completableOptions, ok := a.options.(CompletableOptions); ok {
		if err := completableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) > 0 {
		var aggerr error
		for _, err := range errs {
			aggerr = multierr.Append(aggerr, err)
		}
		return aggerr
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
