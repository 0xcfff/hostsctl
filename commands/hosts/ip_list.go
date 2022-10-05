package hosts

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/0xcfff/dnspipe/backend/hosts"
	common "github.com/0xcfff/dnspipe/commands/cmdcommon"
	"github.com/0xcfff/dnspipe/printutil"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type outFormat int

const (
	fmtText  outFormat = iota
	fmtHosts outFormat = iota
	fmtJson  outFormat = iota
	fmtYaml  outFormat = iota
)

type outGrouping int

const (
	grpOrig    outGrouping = iota
	grpUngroup outGrouping = iota
	grpGroup   outGrouping = iota
)

var (
	formats = map[string]outFormat{
		"":              fmtText,
		common.TfmtText: fmtText,
		common.TfmtJson: fmtJson,
		common.TfmtYaml: fmtYaml,
		"hosts":         fmtHosts,
	}

	groupings = map[string]outGrouping{
		"":        grpOrig,
		"orig":    grpOrig,
		"group":   grpGroup,
		"ungrpup": grpUngroup,
	}
)

type IpListOptions struct {
	outputFormat   string
	output         outFormat
	outputGrouping string
	grouping       outGrouping
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

	err = printutil.PrintTabbed(os.Stdout, nil, 2, func(w io.Writer) error {
		columns := []string{"IP", "HOSTNAME", "SOURCE", "COMMENT"}
		fmt.Fprint(w, strings.Join(columns, "\t"))
		fmt.Fprintln(w)

		allIps := slices.Clone(c.IPRecords)
		sort.Slice(allIps, func(i, j int) bool { return allIps[i].IP < allIps[j].IP })

		for _, ip := range c.IPRecords {
			aliases := slices.Clone(ip.Aliases)
			sort.Strings(aliases)

			for _, al := range aliases {
				values := []string{ip.IP, al, ip.Notes}
				fmt.Fprint(w, strings.Join(values, "\t"))
				fmt.Fprintln(w)
			}
		}

		return nil
	})
	cobra.CheckErr(err)

	return nil

	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	// b := strings.Builder{}
	// w.Write([]byte("IP\tHostname\tAliases\tNotes\n"))
	// for _, ip := range c.IPRecords {
	// 	b.Reset()
	// 	b.WriteString(ip.IP)
	// 	b.WriteRune('\t')
	// 	if len(ip.Aliases) > 0 {
	// 		b.WriteString(ip.Aliases[0])
	// 	}
	// 	b.WriteRune('\t')
	// 	if len(ip.Aliases) > 1 {
	// 		a := ip.Aliases[1:]
	// 		first := true
	// 		for _, al := range a {
	// 			if !first {
	// 				b.WriteRune(' ')
	// 			} else {
	// 				first = false
	// 			}
	// 			b.WriteString(al)
	// 		}
	// 	}
	// 	b.WriteRune('\t')
	// 	if ip.Notes != "" {
	// 		b.WriteString(ip.Notes)
	// 	}
	// 	b.WriteRune('\t')
	// 	b.WriteRune('\n')
	// 	w.Write([]byte(b.String()))
	// }

	// widths := w.RememberedWidths()
	// w.Flush()

	// for _, src := range c.SyncBlocks {
	// 	if src.Data != nil && len(src.Data.IPRecords) > 0 {
	// 		fmt.Println("records from a sync source")
	// 		w.SetRememberedWidths(widths)

	// 		for _, ip := range src.Data.IPRecords {
	// 			b.Reset()
	// 			b.WriteString(ip.IP)
	// 			b.WriteRune('\t')
	// 			if len(ip.Aliases) > 0 {
	// 				b.WriteString(ip.Aliases[0])
	// 			}
	// 			b.WriteRune('\t')
	// 			if len(ip.Aliases) > 1 {
	// 				a := ip.Aliases[1:]
	// 				first := true
	// 				for _, al := range a {
	// 					if !first {
	// 						b.WriteRune(' ')
	// 					} else {
	// 						first = false
	// 					}
	// 					b.WriteString(al)
	// 				}
	// 			}
	// 			b.WriteRune('\t')
	// 			if ip.Notes != "" {
	// 				b.WriteString(ip.Notes)
	// 			}
	// 			b.WriteRune('\t')
	// 			b.WriteRune('\n')
	// 			w.Write([]byte(b.String()))
	// 		}
	// 		w.Flush()
	// 	}
	// }

}
