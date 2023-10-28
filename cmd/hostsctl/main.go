package main

import (
	"os"

	"github.com/0xcfff/hostsctl/commands"
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
	VERSION string = "developer <no version>"
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func init() {
	rootCmd = commands.NewCmdRoot(commands.RootParams{
		Version: VERSION,
	})
}
