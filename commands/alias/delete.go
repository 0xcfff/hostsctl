package alias

import (
	"fmt"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/cobra"
)

type AliasDeleteOptions struct {
	command       *cobra.Command
	blockIdOrName string
	force         bool
}

func NewCmdAliasDelete() *cobra.Command {

	opt := &AliasDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete [ip or alias]",
		Short: fmt.Sprintf("Removes IP alias from %s file", hosts.EtcHosts.Path()),
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.blockIdOrName, "block", "b", opt.blockIdOrName, "Block id or name")
	cmd.Flags().BoolVarP(&opt.force, "force", "f", opt.force, "Force command to succeed even if IP alias is not found")

	return cmd
}

func (opt *AliasDeleteOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd

	return nil
}

func (opt *AliasDeleteOptions) Validate() error {
	return nil
}

func (opt *AliasDeleteOptions) Execute() error {

	ipOrAlias, err := readIpOrAlias(opt)
	cobra.CheckErr(err)

	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	ipsBlock, err := findOrCreateTargetAliasesBlock(doc, opt.blockIdOrName, false)
	cobra.CheckErr(err)

	// TODO: Implement this
	_ = ipOrAlias
	_ = ipsBlock

	return nil
}

func readIpOrAlias(opt *AliasDeleteOptions) (string, error) {
	args := opt.command.Flags().Args()
	if len(args) < 1 {
		return "", common.ErrIpOrAliasExpected
	}
	if len(args) > 1 {
		return "", common.ErrTooManyArguments
	}

	return args[0], nil
}
