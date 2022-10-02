package hosts

import (
	"fmt"
	"os"
	"strings"

	"github.com/0xcfff/dnspipe/backend/hosts"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCmdIpList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists IP addresses and aliases defined in /etc/hosts",
		Run: func(cmd *cobra.Command, args []string) {

			f, err := afero.NewOsFs().Open(hosts.EtcHostsPath())
			cobra.CheckErr(err)
			c, err := hosts.ParseHostsFileWithSources(f, hosts.Safe)
			cobra.CheckErr(err)
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			b := strings.Builder{}
			w.Write([]byte("IP\tHostname\tAliases\tNotes\n"))
			for _, ip := range c.IPRecords {
				b.Reset()
				b.WriteString(ip.IP)
				b.WriteRune('\t')
				if len(ip.Aliases) > 0 {
					b.WriteString(ip.Aliases[0])
				}
				b.WriteRune('\t')
				if len(ip.Aliases) > 1 {
					a := ip.Aliases[1:]
					first := true
					for _, al := range a {
						if !first {
							b.WriteRune(' ')
						} else {
							first = false
						}
						b.WriteString(al)
					}
				}
				b.WriteRune('\t')
				if ip.Notes != "" {
					b.WriteString(ip.Notes)
				}
				b.WriteRune('\t')
				b.WriteRune('\n')
				w.Write([]byte(b.String()))
			}

			widths := w.RememberedWidths()
			w.Flush()

			for _, src := range c.SyncBlocks {
				if src.Data != nil && len(src.Data.IPRecords) > 0 {
					fmt.Println("records from a sync source")
					w.SetRememberedWidths(widths)

					for _, ip := range src.Data.IPRecords {
						b.Reset()
						b.WriteString(ip.IP)
						b.WriteRune('\t')
						if len(ip.Aliases) > 0 {
							b.WriteString(ip.Aliases[0])
						}
						b.WriteRune('\t')
						if len(ip.Aliases) > 1 {
							a := ip.Aliases[1:]
							first := true
							for _, al := range a {
								if !first {
									b.WriteRune(' ')
								} else {
									first = false
								}
								b.WriteString(al)
							}
						}
						b.WriteRune('\t')
						if ip.Notes != "" {
							b.WriteString(ip.Notes)
						}
						b.WriteRune('\t')
						b.WriteRune('\n')
						w.Write([]byte(b.String()))
					}
					w.Flush()
				}
			}

		},
	}
	return cmd
}
