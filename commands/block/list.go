package block

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iotools"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type outFormat int

const (
	fmtText  outFormat = iota
	fmtShort outFormat = iota
	fmtWide  outFormat = iota
	fmtJson  outFormat = iota
	fmtYaml  outFormat = iota
)

var (
	formats = map[string]outFormat{
		"":              fmtText,
		common.TfmtText: fmtText,
		"short":         fmtShort,
		"wide":          fmtWide,
		common.TfmtJson: fmtJson,
		common.TfmtYaml: fmtYaml,
	}
)

type BlockListOptions struct {
	command      *cobra.Command
	output       string
	outputFormat outFormat
	noHeaders    bool
}

func NewCmdBlockList() *cobra.Command {

	opt := &BlockListOptions{}

	cmd := &cobra.Command{
		Use:   "list [(-o|--output)=name] [filter]",
		Short: fmt.Sprintf("Lists IP aliases blocks defined in %s", hosts.EtcHosts.Path()),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().BoolVar(&opt.noHeaders, "no-headers", opt.noHeaders, "Disable printing headers")
	cmd.Flags().StringVarP(&opt.output, "output", "o", opt.output, fmt.Sprintf("Output format. One of %s", strings.Join(maps.Keys(formats), ",")))

	return cmd
}

func (opt *BlockListOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd

	var ok bool
	opt.outputFormat, ok = formats[opt.output]
	if !ok {
		return fmt.Errorf("value %v is not support; %w", opt.output, common.ErrNotSupportedOutputFormat)
	}

	return nil
}

func (opt *BlockListOptions) Validate() error {
	return nil
}

func (opt *BlockListOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))
	c, err := src.Load()
	cobra.CheckErr(err)

	switch opt.outputFormat {
	case fmtText, fmtShort, fmtWide:
		err = writeDataAsText(opt, c)
	case fmtJson:
		err = writeDataAsJson(opt, c)
	case fmtYaml:
		err = writeDataAsYaml(opt, c)
	default:
		panic("unknown output format")
	}
	cobra.CheckErr(err)

	return nil
}

func writeDataAsText(opt *BlockListOptions, data *dom.Document) error {
	m := NewBlocksModels(data)

	err := iotools.PrintTabbed(opt.command.OutOrStdout(), nil, 2, func(w io.Writer) error {

		if !opt.noHeaders {
			columns := []string{"ID", "SYS", "NAME", "COMMENT", "ALIASES", "SYSTEM ALIASES"}
			visible := getVisibleValues(opt, columns)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		for _, b := range m {
			sys := ""

			if b.AliasesCount == b.SystemAliasesCount && b.AliasesCount > 0 {
				sys = "+"
			} else if b.SystemAliasesCount > 0 {
				sys = "*"
			}

			values := []string{strconv.Itoa(b.ID), sys, b.Name, b.Comment, strconv.Itoa(b.AliasesCount), strconv.Itoa(b.SystemAliasesCount)}

			visible := getVisibleValues(opt, values)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}
		return nil
	})
	return err
}

func writeDataAsJson(opt *BlockListOptions, data *dom.Document) error {
	m := NewBlocksModels(data)
	buff, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.command.OutOrStdout(), string(buff))
	return nil
}

func writeDataAsYaml(opt *BlockListOptions, data *dom.Document) error {
	m := NewBlocksModels(data)
	buff, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.command.OutOrStdout(), string(buff))
	return nil
}

func getVisibleValues(opt *BlockListOptions, values []string) []string {
	// "ID", "SYS", "NAME", "COMMENT", "ALIASES", "SYSTEM ALIASES"
	switch opt.outputFormat {
	case fmtText:
		return values[:3]
	case fmtShort:
		return []string{values[0], values[2]}
	case fmtWide:
		return values
	default:
		panic("unsupported formatting")
	}
}
