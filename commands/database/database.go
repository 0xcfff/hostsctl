package database

import "github.com/spf13/cobra"

func NewCmdDatabase() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database [command]",
		Short:   "Manage IP aliases database",
		Aliases: []string{"db"},
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdDatabasePrint())
	cmd.AddCommand(NewCmdDatabaseFormat())
	cmd.AddCommand(NewCmdDatabaseLocation())

	return cmd
}
