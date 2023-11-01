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
	ipOrAlias     string
	force         bool
}

func NewCmdAliasDelete() *cobra.Command {

	opt := &AliasDeleteOptions{}

	cmd := &cobra.Command{
		Use:     "delete [ip or alias]",
		Short:   fmt.Sprintf("Removes IP alias from %s file", hosts.EtcHosts.Path()),
		Aliases: []string{"remove", "rm"},
		Args:    cobra.ArbitraryArgs,
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

	var err error
	opt.command = cmd

	opt.ipOrAlias, err = readIpOrAliasArg(opt)
	cobra.CheckErr(err)

	return nil
}

func (opt *AliasDeleteOptions) Validate() error {
	return nil
}

func (opt *AliasDeleteOptions) Execute() error {

	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	entriesMap, err := findEntriesToDelete(doc, opt)
	cobra.CheckErr(err)

	err = validateDelete(entriesMap, opt.ipOrAlias, opt.force)
	cobra.CheckErr(err)

	err = performDelete(entriesMap, opt.ipOrAlias)
	cobra.CheckErr(err)

	doc.Normalize()

	err = src.Save(doc, dom.FmtKeep)
	cobra.CheckErr(err)

	return nil
}

func performDelete(foundEntries map[*dom.IPAliasesBlock][]*dom.IPAliasesEntry, ipOrAlias string) error {
	isIp := iptools.IsIP(ipOrAlias)
	for block, entries := range foundEntries {
		for _, entry := range entries {
			aliases := entry.Aliases()
			if isIp || len(aliases) <= 1 {
				block.RemoveEntry(entry)
			} else {
				entry.RemoveAlias(ipOrAlias)
			}
		}
	}
	return nil
}

func findEntriesToDelete(doc *dom.Document, opt *AliasDeleteOptions) (map[*dom.IPAliasesBlock][]*dom.IPAliasesEntry, error) {

	entriesMap := make(map[*dom.IPAliasesBlock][]*dom.IPAliasesEntry)

	if opt.blockIdOrName != "" {
		block := doc.IPsBlockByIdOrName(opt.blockIdOrName)
		if block == nil {
			if !opt.force {
				return nil, fmt.Errorf("blockId: %s; %w", opt.blockIdOrName, common.ErrBlockNotFound)
			}
		} else {
			entries := block.AliasEntriesByIPOrAlias(opt.ipOrAlias)
			if len(entries) == 0 && !opt.force {
				return nil, common.ErrAliasNotFound
			}
			entriesMap[block] = entries
		}
	} else {
		for _, block := range doc.IPBlocks() {
			var ipsEntries []*dom.IPAliasesEntry = block.AliasEntriesByIPOrAlias(opt.ipOrAlias)
			if len(ipsEntries) > 0 {
				entriesMap[block] = ipsEntries
			}
		}
	}
	return entriesMap, nil
}

func validateDelete(foundEntries map[*dom.IPAliasesBlock][]*dom.IPAliasesEntry, ipOrAlias string, forceFlag bool) error {

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

	if entriesCount == 0 && !forceFlag {
		return common.ErrAliasNotFound
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
