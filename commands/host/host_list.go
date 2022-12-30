package host

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iptools"
	"github.com/0xcfff/hostsctl/printutil"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type outFormat int

const (
	fmtText  outFormat = iota
	fmtShort outFormat = iota
	fmtWide  outFormat = iota
	fmtPlain outFormat = iota
	fmtJson  outFormat = iota
	fmtYaml  outFormat = iota
)

var (
	formats = map[string]outFormat{
		"":              fmtText,
		common.TfmtText: fmtText,
		"short":         fmtShort,
		"wide":          fmtWide,
		"plain":         fmtPlain,
		common.TfmtJson: fmtJson,
		common.TfmtYaml: fmtYaml,
	}

	groupings = map[string]IPGrouping{
		"":         GrpUngroup,
		"orig":     GrpOriginal,
		"original": GrpOriginal,
		"group":    GrpGroup,
		"ungroup":  GrpUngroup,
	}
)

type IpListOptions struct {
	command        *cobra.Command
	output         string
	outputFormat   outFormat
	grouping       string
	outputGrouping IPGrouping
	noHeaders      bool
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
	cmd.Flags().StringVarP(&opt.output, "output", "o", opt.output, fmt.Sprintf("Output format. One of %s", strings.Join(maps.Keys(formats), ",")))
	cmd.Flags().StringVarP(&opt.grouping, "grouping", "g", opt.grouping, fmt.Sprintf("IPs grouping. One of %s", strings.Join(maps.Keys(groupings), ",")))

	return cmd
}

func (opt *IpListOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd

	var ok bool
	opt.outputFormat, ok = formats[opt.output]
	if !ok {
		return fmt.Errorf("--output %v of list command does not support specified output format", opt.output)
	}

	opt.outputGrouping, ok = groupings[opt.grouping]
	if !ok {
		return fmt.Errorf("--grouping %v of list command does not support specified value", opt.grouping)
	}

	return nil
}

func (opt *IpListOptions) Validate() error {
	return nil
}

func (opt *IpListOptions) Execute() error {
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
	case fmtPlain:
		err = writeDataAsHosts(opt, c)
	default:
		panic("unknown output format")
	}
	cobra.CheckErr(err)

	return nil
}

func writeDataAsText(opt *IpListOptions, data *dom.Document) error {
	m := NewHostModels(data, opt.outputGrouping)

	err := printutil.PrintTabbed(opt.command.OutOrStdout(), nil, 2, func(w io.Writer) error {

		if !opt.noHeaders {
			// columns := []string{"IP", "HOSTNAME", "GROUP", "COMMENT"}
			columns := []string{"GRP", "SYS", "IP", "ALIAS", "COMMENT", "GROUP", "GROUP COMMENT"}
			visible := getVisibleValues(opt, columns)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		var prev *HostModel
		var grpId int = 0

		for _, ip := range m {
			grp := ""
			if prev == nil || prev.Group.Id != ip.Group.Id {
				grpId = ip.Group.Id
				grp = fmt.Sprintf("[%v]", grpId)
			}

			sys := ""
			cntSystem := 0

			for _, alias := range ip.Hosts {
				if iptools.IsSystemAlias(ip.IP, alias) {
					cntSystem += 1
				}
			}
			if cntSystem == len(ip.Hosts) {
				sys = "+"
			} else if cntSystem > 0 {
				sys = "*"
			}

			gn := ip.Group.Name
			if gn == "" {
				gn = fmt.Sprint(ip.Group.Id)
			}

			values := []string{grp, sys, ip.IP, strings.Join(ip.Hosts, ", "), ip.Comment, gn, ip.Group.Comment}

			visible := getVisibleValues(opt, values)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
			prev = ip
		}

		return nil
	})
	return err
}

func writeDataAsHosts(opt *IpListOptions, data *dom.Document) error {
	m := NewHostModels(data, opt.outputGrouping)

	err := printutil.PrintTabbed(opt.command.OutOrStdout(), nil, 2, func(w io.Writer) error {
		for _, ip := range m {

			values := []string{ip.IP, strings.Join(ip.Hosts, ", ")}
			fmt.Fprint(w, strings.Join(values, "\t"))
			fmt.Fprintln(w)
		}
		return nil
	})
	return err
}

func writeDataAsJson(opt *IpListOptions, data *dom.Document) error {
	m := NewHostModels(data, opt.outputGrouping)
	buff, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.command.OutOrStdout(), string(buff))
	return nil
}

func writeDataAsYaml(opt *IpListOptions, data *dom.Document) error {
	m := NewHostModels(data, opt.outputGrouping)
	buff, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.command.OutOrStdout(), string(buff))
	return nil
}

func getVisibleValues(opt *IpListOptions, values []string) []string {
	// "GRP", "SYS", "IP", "ALIAS", "COMMENT", "GROUP", "GROUP COMMENT"
	switch opt.outputFormat {
	case fmtText:
		return values[:4]
	case fmtShort:
		return values[2:4]
	case fmtWide:
		return values
	default:
		panic("unsupported formatting")
	}
}
