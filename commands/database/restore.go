package database

import (
	"fmt"
	"io"
	"os"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/cobra"
)

type RestoreOptions struct {
	command *cobra.Command
	source  string
}

func NewCmdDatabaseRestore() *cobra.Command {

	opt := &RestoreOptions{}

	cmd := &cobra.Command{
		Use:   "restore [flags]",
		Short: "Restore IP aliases database",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.source, "source", "s", "", "source backup file name")

	return cmd
}

func (opt *RestoreOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd
	return nil
}

func (opt *RestoreOptions) Validate() error {
	args := opt.command.Flags().Args()
	if len(args) > 0 {
		return common.ErrTooManyArguments
	}
	return nil
}

func (opt *RestoreOptions) Execute() error {

	sourcePath := opt.source
	targetPath := hosts.EtcHosts.Path()
	if sourcePath == "" {
		sourcePath = fmt.Sprintf("%s.bak", targetPath)
	}

	fs := common.FileSystem(opt.command.Context())

	sf, err := fs.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sf.Close()

	tf, err := fs.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer tf.Close()

	err = tf.Truncate(0)
	if err != nil {
		return err
	}

	cnt, err := io.Copy(tf, sf)
	if err != nil {
		return err
	}

	fmt.Printf("%d bytes copied \n", cnt)

	return nil
}
