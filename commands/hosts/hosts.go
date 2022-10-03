package hosts

import "github.com/spf13/cobra"

func NewCmdHosts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hosts [command]",
		Short: "manages /etc/hosts file",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}

	cmd.AddCommand(NewCmdIp())

	return cmd
}
