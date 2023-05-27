package user

import (
	"context"
	"fmt"

	"github.com/changaolee/skeleton/internal/pkg/clioptions"
	"github.com/changaolee/skeleton/internal/sktctl/util"
	"github.com/changaolee/skeleton/internal/sktctl/util/templates"
	apiclientv1 "github.com/changaolee/skeleton/pkg/sdk/apiserver/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	getUsageStr = "get USERNAME"
)

type GetOptions struct {
	Name string

	Client apiclientv1.APIV1Interface
	clioptions.IOStreams
}

var (
	getExample = templates.Examples(`
		# Get user foo detail information
		sktctl user get foo`)

	getUsageErrStr = fmt.Sprintf("expected '%s'.\nUSERNAME is required arguments for the get command", getUsageStr)
)

func NewGetOptions(ioStreams clioptions.IOStreams) *GetOptions {
	return &GetOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGet(f util.Factory, ioStreams clioptions.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   getUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Display a user resource.",
		TraverseChildren:      true,
		Long:                  `Display a user resource.`,
		Example:               getExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(f, cmd, args))
			util.CheckErr(o.Validate(cmd, args))
			util.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

func (o *GetOptions) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return util.UsageErrorf(cmd, getUsageErrStr)
	}

	o.Name = args[0]

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = apiclientv1.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

func (o *GetOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *GetOptions) Run(args []string) error {
	user, err := o.Client.Users().Get(context.TODO(), o.Name)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(o.Out)

	data := [][]string{
		{
			user.Name,
			user.Nickname,
			user.Email,
			user.Phone,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	table = setHeader(table)
	table = util.TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}
