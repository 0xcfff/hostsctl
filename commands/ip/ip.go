package ip

import "github.com/spf13/cobra"

type ipCommandContextValue int

const (
	ctxCustomFs ipCommandContextValue = iota
)

func NewCmdIp() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ip [command]",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdIpList())
	cmd.AddCommand(NewCmdIpAdd())

	return cmd
}
