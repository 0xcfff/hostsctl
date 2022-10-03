package main

import (
	"fmt"
	"os"

	hsts "github.com/0xcfff/dnspipe/commands/hosts"
	"github.com/0xcfff/dnspipe/commands/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("error executing command")
		os.Exit(1)
	}
	fmt.Println("executing successfully")
	os.Exit(0)
}

func init() {
	rootCmd = &cobra.Command{
		Short: "syncdns synchronizes dns records from various sources to local stores (/etc/hosts, text files, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(cmd.Help())
		},
	}

	rootCmd.AddCommand(version.NewCmdVersion())
	rootCmd.AddCommand(hsts.NewCmdHosts())
}
