package hosts

import (
	"fmt"

	"github.com/0xcfff/dnspipe/backend/hosts"
	"github.com/spf13/cobra"
)

func NewCmdIpAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [ip] [alias]...",
		Short: "Adds IP address and aliases to /etc/hosts file",
		Run: func(cmd *cobra.Command, args []string) {
			ip := &hosts.IPRecord{
				IP:      args[0],
				Aliases: args[1:],
			}
			f, err := hosts.NewHostsFile("", nil)
			cobra.CheckErr(err)

			f.AppendIp(ip)
			f.Dump()

			cobra.CheckErr(fmt.Errorf("Not Implemented"))
		},
	}
	return cmd
}
