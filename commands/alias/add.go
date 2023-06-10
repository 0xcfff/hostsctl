package alias

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iptools"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type AliasAddOptions struct {
	command       *cobra.Command
	blockIdOrName string
	comment       string
	force         bool
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
	cmd.Flags().BoolVarP(&opt.force, "force", "", opt.force, "Enforces creation of a named IP block if it is missing")

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

	ipsBlock, err := findOrCreateTargetAliasesBlock(doc, opt.blockIdOrName, opt.force)
	if err != nil {
		return err
	}

	aliases, err := readIpAliases(opt)
	if err != nil {
		return err
	}

	for _, a := range aliases {
		ipsBlock.AddEntry(a)
	}

	fmt.Printf("OS Args\n%v \n", os.Args)
	fmt.Printf("CMD Args\n%v \n", opt.command.Flags().Args())
	fi, err := os.Stdin.Stat()
	if err == nil {
		fmt.Printf("Inpput Size: %d\n", fi.Size())
	}
	fmt.Println("-----------")

	dom.Write(os.Stdout, doc, dom.FmtKeep)

	// TODO: add logic to output result

	return nil
}

func readIpAliases(opt *AliasAddOptions) ([]*dom.IPAliasesEntry, error) {
	// try read IP alias from opts
	if args := opt.command.Flags().Args(); len(args) >= 2 {
		alias, err := readIpAliasFromArgs(opt)
		if err != nil {
			return nil, err
		}
		aliases := make([]*dom.IPAliasesEntry, 1)
		aliases[0] = alias
		return aliases, nil
	}
	return readIpAliasesFromPassedInput(opt)
}

func readIpAliasFromArgs(opt *AliasAddOptions) (*dom.IPAliasesEntry, error) {
	args := opt.command.Flags().Args()
	if len(args) < 2 {
		return nil, errors.New("no data provided in args")
	}
	if !iptools.IsIP(args[0]) {
		return nil, fmt.Errorf("%s is not an IP", args[0])
	}
	alias := dom.NewIPAliasesEntry(args[0])
	for _, a := range args[1:] {
		alias.AddAlias(a)
	}
	if note := opt.comment; note != "" {
		alias.SetNote(note)
	}

	return alias, nil
}

func readIpAliasesFromPassedInput(opt *AliasAddOptions) ([]*dom.IPAliasesEntry, error) {
	r := opt.command.InOrStdin()
	doc, err := dom.Read(r)
	if err != nil {
		return nil, err
	}

	aliases := make([]*dom.IPAliasesEntry, 0)

	for _, b := range doc.Blocks() {
		switch b.Type() {
		case dom.IPList:
			ips := b.(*dom.IPAliasesBlock)
			entries := ips.Entries()
			for _, a := range entries {
				a.ClearFormatting()
			}
			aliases = append(aliases, entries...)
		case dom.Comments:
			{
			}
		case dom.Blanks:
			{
			}
		case dom.Unknown:
			if !opt.force {
				unk := b.(*dom.UnrecognizedBlock)
				elements := unk.BodyElements()
				lineNum := -1
				if len(elements) > 0 {
					lineNum = elements[0].OriginalLineIndex()
				}
				return nil, fmt.Errorf("error in input line %d", lineNum)
			}
		default:
			panic("unknown block type")
		}
	}

	if len(aliases) == 0 && !opt.force {
		return nil, errors.New("no ips aliases provided")
	}

	return aliases, nil
}

func findOrCreateTargetAliasesBlock(doc *dom.Document, ipBlockIdOrName string, createNamedIfMissing bool) (*dom.IPAliasesBlock, error) {

	// #1 try to find ips block by id
	var ipsBlock *dom.IPAliasesBlock
	if ipBlockIdOrName != "" {
		ipsBlock = doc.IPsBlockByIdOrName(ipBlockIdOrName)
		if ipsBlock == nil && !createNamedIfMissing {
			return nil, fmt.Errorf("aliases block '%s' was not found", ipBlockIdOrName)
		}
	}

	// #2 try to find last ips block
	if ipsBlock == nil {
		blocks := doc.IPBlocks()

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

	}

	// #3 try to create a new block
	if ipsBlock == nil {
		ipsBlock := dom.NewIPAliasesBlock()
		if ipBlockIdOrName != "" {
			v, err := strconv.Atoi(ipBlockIdOrName)
			if err != nil {
				ipsBlock.SetId(v)
			} else {
				ipsBlock.SetName(ipBlockIdOrName)
			}
		}
		doc.AddBlock(ipsBlock)
	}
	return ipsBlock, nil
}
