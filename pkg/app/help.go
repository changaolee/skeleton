package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagHelp          = "help"
	flagHelpShorthand = "H"
)

func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command.",
		Long: `Help provides help for any command in the application.
Simply type ` + name + ` help [path to command] for full details.`,

		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				_ = c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag()
				_ = cmd.Help()
			}
		},
	}
}

// addHelpFlag 为指定应用程序添加 help 标志.
func addHelpFlag(name string, fs *pflag.FlagSet) {
	fs.BoolP(flagHelp, flagHelpShorthand, false, fmt.Sprintf("Help for %s.", name))
}

// addHelpCommandFlag 为应用程序的指定命令添加 help 标志.
func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(flagHelp, flagHelpShorthand, false,
		fmt.Sprintf("Help for the %s command.", color.GreenString(strings.Split(usage, " ")[0])),
	)
}

// TerminalSize 返回用户终端的当前宽度和高度.
func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}
