package database

import (
	"fmt"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/cobra"
)

type LocationOptions struct {
	command *cobra.Command
}

func NewCmdDatabaseLocation() *cobra.Command {

	opt := &LocationOptions{}

	cmd := &cobra.Command{
		Use:   "location",
		Short: "Prints IP aliases database location",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	return cmd
}

func (opt *LocationOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd
	return nil
}

func (opt *LocationOptions) Validate() error {
	args := opt.command.Flags().Args()
	if len(args) > 0 {
		return common.ErrTooManyArguments
	}
	return nil
}

func (opt *LocationOptions) Execute() error {

	out := opt.command.OutOrStdout()
	fmt.Fprintln(out, hosts.EtcHosts.Path())

	return nil
}
