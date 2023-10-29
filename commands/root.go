package commands

import (
	"github.com/0xcfff/hostsctl/commands/alias"
	"github.com/0xcfff/hostsctl/commands/block"
	"github.com/0xcfff/hostsctl/commands/database"
	"github.com/0xcfff/hostsctl/commands/version"
	"github.com/spf13/cobra"
)

type RootParams struct {
	Version string
}

func NewCmdRoot(p RootParams) *cobra.Command {
	cmd := &cobra.Command{
		Short: "hostsctl manages ip to hostname mappings (usually stored in /etc/hosts)",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}

	cmd.AddCommand(version.NewCmdVersion(version.VersionParams{
		Version: p.Version,
	}))
	cmd.AddCommand(block.NewCmdBlock())
	cmd.AddCommand(alias.NewCmdAlias())
	cmd.AddCommand(database.NewCmdDatabase())
	return cmd
}
