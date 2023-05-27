package user

import (
	"github.com/changaolee/skeleton/internal/pkg/clioptions"
	"github.com/changaolee/skeleton/internal/sktctl/util"
	"github.com/changaolee/skeleton/internal/sktctl/util/templates"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var userLong = templates.LongDesc(`
	User management commands.

Administrator can use all subcommands, non-administrator only allow to use create/get/update. When call get/update non-administrator only allow to operate their own resources, if permission not allowed, will return an 'Permission denied' error.`)

func NewCmdUser(f util.Factory, ioStreams clioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "user SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage users on skt platform",
		Long:                  userLong,
		Run:                   util.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	cmd.AddCommand(NewCmdCreate(f, ioStreams))
	cmd.AddCommand(NewCmdGet(f, ioStreams))
	//cmd.AddCommand(NewCmdList(f, ioStreams))
	//cmd.AddCommand(NewCmdDelete(f, ioStreams))
	//cmd.AddCommand(NewCmdUpdate(f, ioStreams))

	return cmd
}

// setHeader set headers for user commands.
func setHeader(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{"Name", "Nickname", "Email", "Phone", "Created", "Updated"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgWhiteColor})

	return table
}
