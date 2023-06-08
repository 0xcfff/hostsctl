package host

import (
	"fmt"
	"os"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type AliasAddOptions struct {
	command       *cobra.Command
	blockIdOrName string
	comment       string
}

func NewCmdAliasAdd() *cobra.Command {

	opt := &AliasAddOptions{}

	cmd := &cobra.Command{
		Use:   "add [ip] [alias, ...]",
		Short: fmt.Sprintf("Adds IP alias to %s file", hosts.EtcHosts.Path()),
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.blockIdOrName, "block", "b", opt.blockIdOrName, "Block id or name")
	cmd.Flags().StringVarP(&opt.comment, "comment", "c", opt.comment, "Alias comment")

	return cmd
}

func (opt *AliasAddOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd

	return nil
}

func (opt *AliasAddOptions) Validate() error {
	return nil
}

func (opt *AliasAddOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	doc, err := src.Load()
	cobra.CheckErr(err)

	var ipsBlock *dom.IPAliasesBlock
	if opt.blockIdOrName != "" {
		ipsBlock = doc.IPsBlockByIdOrName(opt.blockIdOrName)
		if ipsBlock == nil {
			return fmt.Errorf("block '%s' was not found", opt.blockIdOrName)
		}
	} else {
		blocks := doc.IPBlocks()

		// try to use last IPS bock
		if len(blocks) > 0 {
			lastIPsBlock := blocks[len(blocks)-1]
			var lastIPBlockIfx dom.Block = lastIPsBlock

			allBlocks := doc.Blocks()
			lastIndex := slices.Index(allBlocks, lastIPBlockIfx)
			foundBreakingBlock := false
			for i := lastIndex + 1; i < len(allBlocks); i++ {
				blockType := allBlocks[i].Type()
				shouldBrak := false
				switch blockType {
				case dom.Blanks:
					continue
				default:
					foundBreakingBlock = true
					shouldBrak = true
				}
				if shouldBrak {
					break
				}
			}
			if !foundBreakingBlock {
				ipsBlock = lastIPsBlock
			}
		}

		// create new block
		if ipsBlock == nil {
			ipsBlock := dom.NewIPAliasesBlock()
			doc.AddBlock(ipsBlock)
		}
	}

	fmt.Printf("OS Args\n%v \n", os.Args)
	fmt.Printf("CMD Args\n%v \n", opt.command.Flags().Args())

	ipAlias := dom.NewIPAliasesEntry("127.0.0.1")
	ipAlias.AddAlias("test1.local")
	ipAlias.AddAlias("test2.local")
	ipsBlock.AddEntry(ipAlias)

	dom.Write(os.Stdout, doc, dom.FmtDefault)

	// TODO: add logic to output result

	return nil
}
