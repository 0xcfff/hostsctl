package commands

import (
	"github.com/0xcfff/hostsctl/commands/alias"
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
	cmd.AddCommand(alias.NewCmdAlias())
	return cmd
}
