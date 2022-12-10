package main

import (
	"os"

	"github.com/0xcfff/hostsctl/commands/ip"
	"github.com/0xcfff/hostsctl/commands/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func init() {
	rootCmd = &cobra.Command{
		Short: "hostsctl manages ip to hostname mappings (usually stored in /etc/hosts)",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}

	rootCmd.AddCommand(version.NewCmdVersion())
	rootCmd.AddCommand(ip.NewCmdIp())
}
