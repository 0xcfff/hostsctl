package common

import (
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	OutputFormat string
}

type CliCommand interface {
	Complete(cmd *cobra.Command, args []string) error
	Validate() error
	Execute() error
}

func RunCliCommand(cliCmd CliCommand, cmd *cobra.Command, args []string) {
	cobra.CheckErr(cliCmd.Complete(cmd, args))
	cobra.CheckErr(cliCmd.Validate())
	cobra.CheckErr(cliCmd.Execute())
}
