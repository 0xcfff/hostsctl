package host

import (
	"github.com/spf13/cobra"
)

type AliasAddOptions struct {
	command        *cobra.Command
	output         string
	outputFormat   outFormat
	grouping       string
	outputGrouping IPGrouping
	noHeaders      bool
}

func NewCmdAliasAdd() *cobra.Command {

	opt := &AliasAddOptions{}

	cmd := &cobra.Command{
		// TODO: review below
		Use:   "add [ip] [alias]...",
		Short: "Adds IP address and aliases to /etc/hosts file",
		Run: func(cmd *cobra.Command, args []string) {
			// ip := &hosts.IPRecord{
			// 	IP:      args[0],
			// 	Aliases: args[1:],
			// }
			// fs := hosts.NewHostsFileSource("", nil)
			// f, err := fs.LoadFile()
			// cobra.CheckErr(err)

			// f.AppendIp(ip)
			// f.Dump()

			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}
	return cmd
}

func (opt *AliasAddOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (opt *AliasAddOptions) Validate() error {
	return nil
}

func (opt *AliasAddOptions) Execute() error {
	return nil
}
