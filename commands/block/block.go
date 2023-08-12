package block

import "github.com/spf13/cobra"

func NewCmdBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [command]",
		Short: "Manage IP aliases blocks",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdBlockList())
	cmd.AddCommand(NewCmdBlockAdd())
	cmd.AddCommand(NewCmdBlockDelete())

	return cmd
}
