package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/printutil"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type outFormat int

const (
	fmtText  outFormat = iota
	fmtHosts outFormat = iota
	fmtJson  outFormat = iota
	fmtYaml  outFormat = iota
)

var (
	formats = map[string]outFormat{
		"":              fmtText,
		common.TfmtText: fmtText,
		common.TfmtJson: fmtJson,
		common.TfmtYaml: fmtYaml,
		"hosts":         fmtHosts,
	}

	groupings = map[string]IPGrouping{
		"":         GrpOriginal,
		"orig":     GrpOriginal,
		"original": GrpOriginal,
		"combine":  GrpGroup,
		"ungroup":  GrpUngroup,
	}
)

type IpListOptions struct {
	cmd            *cobra.Command
	outputFormat   string
	output         outFormat
	outputGrouping string
	grouping       IPGrouping
	noHeaders      bool
	noGroup        bool
	noComments     bool
}

func NewCmdIpList() *cobra.Command {

	opt := &IpListOptions{}

	cmd := &cobra.Command{
		Use:   "list [(-o|--output)=name] [filter]",
		Short: fmt.Sprintf("Lists IP addresses and aliases defined in %s", hosts.EtcHosts.Path()),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().BoolVar(&opt.noHeaders, "no-headers", opt.noHeaders, "Disable printing headers")
	cmd.Flags().BoolVar(&opt.noGroup, "no-group", opt.noGroup, "Do not show IP group")
	cmd.Flags().BoolVar(&opt.noComments, "no-comments", opt.noComments, "Do not show comments")
	cmd.Flags().StringVarP(&opt.outputFormat, "output", "o", opt.outputFormat, fmt.Sprintf("Output format. One of %s", strings.Join(maps.Keys(formats), ",")))
	cmd.Flags().StringVarP(&opt.outputGrouping, "grouping", "g", opt.outputFormat, fmt.Sprintf("IPs grouping. One of %s", strings.Join(maps.Keys(groupings), ",")))

	return cmd
}

func (opt *IpListOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.cmd = cmd

	var ok bool
	opt.output, ok = formats[opt.outputFormat]
	if !ok {
		return fmt.Errorf("--output %v of list command does not support specified output format", opt.outputFormat)
	}

	opt.grouping, ok = groupings[opt.outputGrouping]
	if !ok {
		return fmt.Errorf("--grouping %v of list command does not support specified value", opt.outputGrouping)
	}

	return nil
}

func (opt *IpListOptions) Validate() error {
	return nil
}

func (opt *IpListOptions) Execute() error {
	c, err := hosts.EtcHosts.Load()
	cobra.CheckErr(err)

	switch opt.output {
	case fmtText:
		err = writeDataAsText(opt, c)
	case fmtJson:
		err = writeDataAsJson(opt, c)
	case fmtYaml:
		err = writeDataAsYaml(opt, c)
	case fmtHosts:
		err = writeDataAsHosts(opt, c)
	default:
		panic("unknown output format")
	}
	cobra.CheckErr(err)

	return nil

}

func writeDataAsText(opt *IpListOptions, data *dom.Document) error {
	m := NewIPModels(data, opt.grouping)

	err := printutil.PrintTabbed(opt.cmd.OutOrStdout(), nil, 2, func(w io.Writer) error {

		if !opt.noHeaders {
			columns := []string{"IP", "HOSTNAME", "GROUP", "COMMENT"}
			visible := getVisibleValues(opt, columns)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		for _, ip := range m {
			values := []string{ip.IP, strings.Join(ip.Aliases, ", "), ip.Group, ip.Comment}
			visible := getVisibleValues(opt, values)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		return nil
	})
	return err
}

func writeDataAsHosts(opt *IpListOptions, data *dom.Document) error {
	m := NewIPModels(data, opt.grouping)

	panic("not implemented")
	fmt.Println(m)
	return nil
}

func writeDataAsJson(opt *IpListOptions, data *dom.Document) error {
	m := NewIPModels(data, opt.grouping)
	buff, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.cmd.OutOrStdout(), string(buff))
	return nil
}

func writeDataAsYaml(opt *IpListOptions, data *dom.Document) error {
	m := NewIPModels(data, opt.grouping)
	buff, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.cmd.OutOrStdout(), string(buff))
	return nil
}

func getVisibleValues(opt *IpListOptions, values []string) []string {
	// "IP", "HOSTNAME", "SOURCE", "COMMENT"
	if opt.noComments && opt.noGroup {
		return values[:1]
	}
	if opt.noComments {
		return values[:2]
	}
	if opt.noGroup {
		return []string{values[0], values[1], values[3]}
	}
	return values
}
