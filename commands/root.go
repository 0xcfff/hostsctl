package commands

import (
	"github.com/0xcfff/hostsctl/commands/alias"
	"github.com/0xcfff/hostsctl/commands/block"
	"github.com/0xcfff/hostsctl/commands/format"
	"github.com/0xcfff/hostsctl/commands/print"
	"github.com/0xcfff/hostsctl/commands/version"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Short: "hostsctl manages ip to hostname mappings (usually stored in /etc/hosts)",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}

	cmd.AddCommand(version.NewCmdVersion())
	cmd.AddCommand(block.NewCmdBlock())
	cmd.AddCommand(alias.NewCmdAlias())
	cmd.AddCommand(format.NewCmdFormatDocument())
	cmd.AddCommand(print.NewCmdPrintDocument())
	return cmd
}
