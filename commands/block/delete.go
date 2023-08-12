package block

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type BlockDeleteOptions struct {
	command   *cobra.Command
	blockId   int
	blockName string
	force     bool
}

func NewCmdBlockDelete() *cobra.Command {

	opt := &BlockDeleteOptions{}
	opt.blockId = emptyId

	cmd := &cobra.Command{
		Use:   "delete [id or name]",
		Short: fmt.Sprintf("Removes IP aliases block from %s file", hosts.EtcHosts.Path()),
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().BoolVarP(&opt.force, "force", "f", opt.force, "Force command to succeed even if IP alias is not found")

	return cmd
}

func (opt *BlockDeleteOptions) Complete(cmd *cobra.Command, args []string) error {

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

func (opt *BlockDeleteOptions) Validate() error {
	return nil
}

func (opt *BlockDeleteOptions) Execute() error {

	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	targerBlocks, err := findTargetBlockForDelete(doc, opt)
	cobra.CheckErr(err)

	err = deleteBlocks(doc, targerBlocks)
	cobra.CheckErr(err)

	err = src.Save(doc, dom.FmtKeep)
	cobra.CheckErr(err)

	return nil
}

func findTargetBlockForDelete(doc *dom.Document, opt *BlockDeleteOptions) ([]*dom.IPAliasesBlock, error) {
	selectedBlocks := doc.IPBlocksByIdentifiers(opt.blockId, opt.blockName)

	blocksFound := len(selectedBlocks)
	if blocksFound == 0 {
		if !opt.force {
			return nil, common.ErrBlockNotFound
		}
	} else if blocksFound > 1 {
		if !opt.force {
			return nil, fmt.Errorf("%d blocks found matching parameters: %w", blocksFound, common.ErrTooManyEntries)
		}
	} else {
		entriesFound := len(selectedBlocks[0].AliasEntries())
		if entriesFound > 0 && !opt.force {
			return nil, fmt.Errorf("target block is not empty, %d entry(es) found in the block: %w", blocksFound, common.ErrTooManyEntries)
		}
	}

	return selectedBlocks, nil
}

func deleteBlocks(doc *dom.Document, targerBlocks []*dom.IPAliasesBlock) error {
	for _, b := range targerBlocks {
		blocks := doc.Blocks()

		var mainBlock, spaceBlock dom.Block

		mainBlock = b
		idx := slices.Index(blocks, mainBlock)
		if len(blocks) > idx+1 {
			spaceBlock = blocks[idx+1]
		}

		doc.DeleteBlock(mainBlock)
		if spaceBlock != nil && spaceBlock.Type() == dom.Blanks {
			doc.DeleteBlock(spaceBlock)
		}
	}

	return nil
}
