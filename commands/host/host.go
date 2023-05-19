package host

import "github.com/spf13/cobra"

func NewCmdHost() *cobra.Command {
	cmd := &cobra.Command{
		Use: "host [command]",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdIpList())
	cmd.AddCommand(NewCmdAliasAdd())

	return cmd
}
