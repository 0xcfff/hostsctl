package block

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/spf13/cobra"
)

const (
	emptyId = -1
)

type BlockAddOptions struct {
	command   *cobra.Command
	blockId   int
	blockName string
	comment   string
	force     bool
}

func NewCmdBlockAdd() *cobra.Command {

	opt := &BlockAddOptions{}
	opt.blockId = emptyId

	cmd := &cobra.Command{
		Use:   "add [id or name]",
		Short: fmt.Sprintf("Adds IP aliases block to %s file", hosts.EtcHosts.Path()),
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.blockName, "name", "t", opt.blockName, "Block name")
	cmd.Flags().IntVarP(&opt.blockId, "id", "n", opt.blockId, "Block id")
	cmd.Flags().StringVarP(&opt.comment, "comment", "c", opt.comment, "Block comment")
	cmd.Flags().BoolVarP(&opt.force, "force", "f", opt.force, "Do not fail if the block already exists, just update it with provided data")

	return cmd
}

func (opt *BlockAddOptions) Complete(cmd *cobra.Command, args []string) error {

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

func (opt *BlockAddOptions) Validate() error {
	return nil
}

func (opt *BlockAddOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	block, err := findTargetBlockForAdd(doc, opt)
	cobra.CheckErr(err)

	if block != nil {
		if !opt.force {
			return common.ErrEntryAlreadyExists
		}
		err = updateBlock(block, opt)
		cobra.CheckErr(err)
	} else {
		block, err = createBlock(opt)
		cobra.CheckErr(err)

		doc.AddBlock(block)
	}

	doc.Normalize()

	err = src.Save(doc, dom.FmtKeep)
	cobra.CheckErr(err)

	return nil
}

func updateBlock(block *dom.IPAliasesBlock, opts *BlockAddOptions) error {
	if opts.blockId != emptyId {
		block.SetId(opts.blockId)
	}
	if opts.blockName != "" {
		block.SetName(opts.blockName)
	}
	if opts.comment != "" {
		block.SetNote(opts.comment)
	}
	return nil
}

func createBlock(opts *BlockAddOptions) (*dom.IPAliasesBlock, error) {
	block := dom.NewIPAliasesBlock()
	err := updateBlock(block, opts)
	return block, err
}

func findTargetBlockForAdd(doc *dom.Document, opt *BlockAddOptions) (*dom.IPAliasesBlock, error) {
	selectedBlocks := doc.IPBlocksByIdentifiers(opt.blockId, opt.blockName)

	blocksFound := len(selectedBlocks)
	if blocksFound > 1 {
		return nil, fmt.Errorf("multiple blocks found matching criteria: %w", common.ErrTooManyEntries)
	} else if blocksFound == 1 {
		return selectedBlocks[0], nil
	}

	return nil, nil
}
