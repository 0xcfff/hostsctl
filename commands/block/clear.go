package block

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iptools"
	"github.com/spf13/cobra"
)

type BlockClearOptions struct {
	command   *cobra.Command
	blockId   int
	blockName string
	force     bool
}

func NewCmdBlockClear() *cobra.Command {

	opt := &BlockClearOptions{}
	opt.blockId = emptyId

	cmd := &cobra.Command{
		Use:   "clear [id or name]",
		Short: fmt.Sprintf("Clears IP aliases block in %s file", hosts.EtcHosts.Path()),
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.blockName, "name", "t", opt.blockName, "Block name")
	cmd.Flags().IntVarP(&opt.blockId, "id", "n", opt.blockId, "Block id")
	cmd.Flags().BoolVarP(&opt.force, "force", "f", opt.force, "Clear the block even if it has system aliases")

	return cmd
}

func (opt *BlockClearOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd

	parsedArgs := cmd.Flags().Args()
	if len(args) > 1 {
		return common.ErrTooManyArguments
	}

	if len(parsedArgs) == 1 {
		blockIdOrName := parsedArgs[0]
		if id, err := strconv.Atoi(blockIdOrName); err == nil {
			if opt.blockId != emptyId {
				return errors.New("block Id is provided twice")
			}
			opt.blockId = id
		} else {
			if opt.blockName != "" {
				return errors.New("block Name is provided twice")
			}
			opt.blockName = blockIdOrName
		}
	}

	return nil
}

func (opt *BlockClearOptions) Validate() error {
	return nil
}

func (opt *BlockClearOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	block, err := findTargetBlockForClear(doc, opt)
	cobra.CheckErr(err)

	if block == nil {
		if !opt.force {
			return common.ErrBlockNotFound
		}
	} else {
		err = validateClear(block, opt)
		cobra.CheckErr(err)

		err = clearBlock(block, opt)
		cobra.CheckErr(err)
	}

	doc.Normalize()

	err = src.Save(doc, dom.FmtKeep)
	cobra.CheckErr(err)

	return nil
}

func validateClear(block *dom.IPAliasesBlock, opts *BlockClearOptions) error {
	if opts.force {
		return nil
	}

	for _, ent := range block.AliasEntries() {
		ip := ent.IP()
		for _, alias := range ent.Aliases() {
			if iptools.IsSystemAlias(ip, alias) {
				return errors.New("the block has system aliases")
			}
		}
	}
	return nil
}
func clearBlock(block *dom.IPAliasesBlock, opts *BlockClearOptions) error {
	for _, ent := range block.AliasEntries() {
		block.RemoveEntry(ent)
	}
	return nil
}

func findTargetBlockForClear(doc *dom.Document, opt *BlockClearOptions) (*dom.IPAliasesBlock, error) {
	selectedBlocks := doc.IPBlocksByIdentifiers(opt.blockId, opt.blockName)

	blocksFound := len(selectedBlocks)
	if blocksFound > 1 {
		return nil, fmt.Errorf("multiple blocks found matching criteria: %w", common.ErrTooManyEntries)
	} else if blocksFound == 1 {
		return selectedBlocks[0], nil
	}

	return nil, nil
}
