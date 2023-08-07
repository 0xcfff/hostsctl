package alias

import "github.com/spf13/cobra"

func NewCmdAlias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias [command]",
		Short: "Manage IP aliases",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdAliasList())
	cmd.AddCommand(NewCmdAliasAdd())
	cmd.AddCommand(NewCmdAliasDelete())

	return cmd
}
