package alias

import (
	"fmt"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iptools"
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

	ipOrAlias, err := readIpOrAliasArg(opt)
	cobra.CheckErr(err)

	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	entriesMap := make(map[int][]*dom.IPAliasesEntry)
	for idx, blk := range doc.IPBlocks() {
		var ipsEntries []*dom.IPAliasesEntry = blk.EntriesByIPOrAlias(ipOrAlias)
		if len(ipsEntries) > 0 {
			entriesMap[idx] = ipsEntries
		}
	}

	err = validateDelete(entriesMap, ipOrAlias, opt.force)
	cobra.CheckErr(err)

	// TODO: Implement this
	return nil
}

func validateDelete(foundEntries map[int][]*dom.IPAliasesEntry, ipOrAlias string, forceFlag bool) error {

	isAlias := !iptools.IsIP(ipOrAlias)

	entriesCount := 0
	systemCount := 0

	for _, entries := range foundEntries {
		entriesCount += len(entries)

		for _, ipe := range entries {
			if isAlias {
				if iptools.IsSystemAlias(ipe.IP(), ipOrAlias) {
					systemCount += 1
				}
			} else {
				for _, alias := range ipe.Aliases() {
					if iptools.IsSystemAlias(ipOrAlias, alias) {
						systemCount += 1
						break
					}
				}
			}
		}
	}

	if systemCount > 0 && !forceFlag {
		return fmt.Errorf("%d of %d entries is system", systemCount, entriesCount)
	}

	if entriesCount > 1 && !forceFlag {
		return fmt.Errorf("%d entries found; %w", entriesCount, common.ErrTooManyEntries)
	}

	return nil
}

func readIpOrAliasArg(opt *AliasDeleteOptions) (string, error) {
	args := opt.command.Flags().Args()
	if len(args) < 1 {
		return "", common.ErrIpOrAliasExpected
	}
	if len(args) > 1 {
		return "", common.ErrTooManyArguments
	}

	return args[0], nil
}
