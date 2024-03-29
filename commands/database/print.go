package database

import (
	"bufio"
	"fmt"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type PrintOptions struct {
	command *cobra.Command
}

func NewCmdDatabasePrint() *cobra.Command {

	opt := &PrintOptions{}

	cmd := &cobra.Command{
		Use:   "print",
		Short: "Prints contents of the database",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	return cmd
}

func (opt *PrintOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd
	return nil
}

func (opt *PrintOptions) Validate() error {
	args := opt.command.Flags().Args()
	if len(args) > 0 {
		return common.ErrTooManyArguments
	}
	return nil
}

func (opt *PrintOptions) Execute() error {
	src := hosts.NewSource(hosts.EtcHosts.Path(), common.FileSystem(opt.command.Context()))

	out := opt.command.OutOrStdout()

	err := src.Apply(func(path string, fs afero.Fs) error {
		f, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		s := bufio.NewScanner(f)
		for s.Scan() {
			fmt.Fprintln(out, s.Text())
		}
		return nil
	})

	cobra.CheckErr(err)

	return nil
}
