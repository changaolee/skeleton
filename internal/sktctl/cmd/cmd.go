// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/changaolee/skeleton/internal/pkg/clioptions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	genericapiserver "github.com/changaolee/skeleton/internal/pkg/server"
	"github.com/changaolee/skeleton/internal/sktctl/cmd/info"
	"github.com/changaolee/skeleton/internal/sktctl/util"
	"github.com/changaolee/skeleton/pkg/app"
)

const commandDesc = `The sktctl controls the skt platform, is the client side tool for skt platform.`

// NewDefaultSKTCtlCommand 使用默认参数创建一个 `sktctl` 命令.
func NewDefaultSKTCtlCommand() *cobra.Command {
	return NewSKTCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewSKTCtlCommand 返回一个 `sktctl` 命令实例.
func NewSKTCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "sktctl",
		Short: "sktctl controls the skt platform",
		Long:  commandDesc,
		Run:   runHelp,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return util.InitProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return util.FlushProfiling()
		},
	}
	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(app.WarnWordSepNormalizeFunc)

	util.AddProfilingFlags(flags)

	sktConfigFlags := clioptions.NewConfigFlags(true)
	sktConfigFlags.WithDeprecatedPasswordFlag()
	sktConfigFlags.WithDeprecatedSecretFlag()
	sktConfigFlags.AddFlags(flags)

	matchVersionSKTConfigFlags := util.NewMatchVersionFlags(sktConfigFlags)
	matchVersionSKTConfigFlags.AddFlags(cmds.PersistentFlags())

	_ = viper.BindPFlags(cmds.PersistentFlags())
	cobra.OnInitialize(func() {
		genericapiserver.LoadConfig(viper.GetString(clioptions.FlagSKTConfig), "sktctl")
	})
	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	cmds.SetGlobalNormalizationFunc(app.WarnWordSepNormalizeFunc)

	f := util.NewFactory(matchVersionSKTConfigFlags)
	ioStreams := clioptions.IOStreams{
		In:     in,
		Out:    out,
		ErrOut: err,
	}

	groups := util.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				info.NewCmdInfo(f, ioStreams),
			},
		},
	}
	groups.AddTo(cmds)

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
