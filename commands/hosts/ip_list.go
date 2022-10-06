package hosts

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/0xcfff/dnspipe/backend/hosts"
	"github.com/0xcfff/dnspipe/commands/common"
	"github.com/0xcfff/dnspipe/printutil"
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
		"":        GrpOriginal,
		"orig":    GrpOriginal,
		"group":   GrpGroup,
		"ungroup": GrpUngroup,
	}
)

type IpListOptions struct {
	outputFormat   string
	output         outFormat
	outputGrouping string
	grouping       IPGrouping
	noHeaders      bool
	noSource       bool
	noComments     bool
}

func NewCmdIpList() *cobra.Command {

	opt := &IpListOptions{}

	cmd := &cobra.Command{
		Use:   "list [(-o|--output)=name] [filter]",
		Short: "Lists IP addresses and aliases defined in /etc/hosts",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().BoolVar(&opt.noHeaders, "no-headers", opt.noHeaders, "Disable printing headers")
	cmd.Flags().BoolVar(&opt.noSource, "no-source", opt.noSource, "Do not show IP sources")
	cmd.Flags().BoolVar(&opt.noComments, "no-comments", opt.noSource, "Do not show comments")
	cmd.Flags().StringVarP(&opt.outputFormat, "output", "o", opt.outputFormat, fmt.Sprintf("Output format. One of %s", strings.Join(maps.Keys(formats), ",")))
	cmd.Flags().StringVarP(&opt.outputGrouping, "grouping", "g", opt.outputFormat, fmt.Sprintf("IPs grouping. One of %s", strings.Join(maps.Keys(groupings), ",")))

	return cmd
}

func (opt *IpListOptions) Complete(cmd *cobra.Command, args []string) error {

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
	fs := hosts.NewHostsFileSource("", nil)
	// f, err := fs.LoadFile()
	// cobra.CheckErr(err)
	// cips :=
	c, err := fs.LoadAndParse(hosts.Strict, hosts.None)
	cobra.CheckErr(err)

	switch opt.output {
	case fmtText:
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

func writeDataAsText(opt *IpListOptions, data *hosts.HostsFileContent) error {
	m := NewIPModels(data, opt.grouping)

	err := printutil.PrintTabbed(os.Stdout, nil, 2, func(w io.Writer) error {

		if !opt.noHeaders {
			columns := []string{"IP", "HOSTNAME", "SOURCE", "COMMENT"}
			visible := getVisibleValues(opt, columns)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		for _, ip := range m {
			values := []string{ip.IP, strings.Join(ip.Aliases, ", "), ip.Source, ip.Comment}
			visible := getVisibleValues(opt, values)
			fmt.Fprint(w, strings.Join(visible, "\t"))
			fmt.Fprintln(w)
		}

		return nil
	})
	return err
}

func writeDataAsJson(opt *IpListOptions, data *hosts.HostsFileContent) error {
	m := NewIPModels(data, opt.grouping)
	buff, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(buff))
	return nil
}

func writeDataAsYaml(opt *IpListOptions, data *hosts.HostsFileContent) error {
	m := NewIPModels(data, opt.grouping)
	buff, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Println(string(buff))
	return nil
}

func getVisibleValues(opt *IpListOptions, values []string) []string {
	// "IP", "HOSTNAME", "SOURCE", "COMMENT"
	if opt.noComments && opt.noSource {
		return values[:1]
	}
	if opt.noComments {
		return values[:2]
	}
	if opt.noSource {
		return []string{values[0], values[1], values[3]}
	}
	return values
}
