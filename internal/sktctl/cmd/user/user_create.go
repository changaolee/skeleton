package user

import (
	"context"
	"fmt"

	"github.com/changaolee/skeleton/internal/pkg/clioptions"
	"github.com/changaolee/skeleton/internal/pkg/model/user"
	"github.com/changaolee/skeleton/internal/pkg/rest"
	"github.com/changaolee/skeleton/internal/sktctl/util"
	"github.com/changaolee/skeleton/internal/sktctl/util/templates"
	metav1 "github.com/changaolee/skeleton/pkg/meta/v1"
	"github.com/spf13/cobra"
)

const (
	createUsageStr = "create USERNAME PASSWORD EMAIL"
)

type CreateOptions struct {
	Nickname string
	Phone    string

	User *user.User

	client rest.Interface
	clioptions.IOStreams
}

var (
	createLong = templates.LongDesc(`Create a user on skt platform.
If nickname not specified, username will be used.`)

	createExample = templates.Examples(`
		# Create user with given input
		sktctl user create foo Foo@2023 foo@test.com

		# Create user wt 
		sktctl user create foo Foo@2023 foo@test.com --phone=18888888xxx --nickname=sktxxx`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERNAME, PASSWORD and EMAIL are required arguments for the create command",
		createUsageStr,
	)
)

func NewCreateOptions(ioStreams clioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdCreate(f util.Factory, ioStreams clioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Create a user resource",
		TraverseChildren:      true,
		Long:                  createLong,
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(f, cmd, args))
			util.CheckErr(o.Validate(cmd, args))
			util.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Nickname, "nickname", o.Nickname, "The nickname of the user.")
	cmd.Flags().StringVar(&o.Phone, "phone", o.Phone, "The phone number of the user.")

	return cmd
}

func (o *CreateOptions) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) < 3 {
		return util.UsageErrorf(cmd, createUsageErrStr)
	}

	if o.Nickname == "" {
		o.Nickname = args[0]
	}

	o.User = &user.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Nickname: o.Nickname,
		Password: args[1],
		Email:    args[2],
		Phone:    o.Phone,
	}

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.client, err = rest.RESTClientFor(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if errs := o.User.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

func (o *CreateOptions) Run(args []string) error {
	var u *user.User

	err := o.client.Post().
		AbsPath("/v1/users").
		Body(o.User).
		Do(context.TODO()).
		Into(&u)

	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(o.Out, "user/%s created\n", u.Name)

	return nil
}
