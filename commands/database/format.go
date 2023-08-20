package database

import (
	"fmt"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/spf13/cobra"
)

type FormatOptions struct {
	command *cobra.Command
	dryRun  bool
}

func NewCmdDatabaseFormat() *cobra.Command {

	opt := &FormatOptions{}

	cmd := &cobra.Command{
		Use:   "format [--dry-run] [filter]",
		Short: fmt.Sprintf("Formats %s", hosts.EtcHosts.Path()),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().BoolVar(&opt.dryRun, "dry-run", opt.dryRun, "Do not store formatting result, instead prints in to output")

	return cmd
}

func (opt *FormatOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd
	return nil
}

func (opt *FormatOptions) Validate() error {
	args := opt.command.Flags().Args()
	if len(args) > 0 {
		return common.ErrTooManyArguments
	}
	return nil
}

func (opt *FormatOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	c, err := src.Load()
	cobra.CheckErr(err)

	if opt.dryRun {
		dom.Write(opt.command.OutOrStdout(), c, dom.FmtReFormat)
	} else {
		src.Save(c, dom.FmtReFormat)
	}

	return nil
}
