package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command 是命令行应用程序中子命令的主体结构.
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

// RunCommandFunc 定义命令执行的回调函数.
type RunCommandFunc func(args []string) error

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false

	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	var namedFlagSets NamedFlagSets
	fs := cmd.Flags()
	if c.options != nil {
		namedFlagSets = c.options.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}
	addHelpCommandFlag(c.usage, fs)

	return cmd
}
