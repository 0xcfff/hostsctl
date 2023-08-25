package database

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type BackupOptions struct {
	command *cobra.Command
	output  string
	force   bool
}

func NewCmdDatabaseBackup() *cobra.Command {

	opt := &BackupOptions{}

	cmd := &cobra.Command{
		Use:   "backup [flags]",
		Short: "Backups IP aliases database",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(opt.Complete(cmd, args))
			cobra.CheckErr(opt.Validate())
			cobra.CheckErr(opt.Execute())
		},
	}

	cmd.Flags().StringVarP(&opt.output, "output", "o", "", "backup file name")
	cmd.Flags().BoolVarP(&opt.force, "force", "f", opt.force, "Do not fail if backup file already exists")

	return cmd
}

func (opt *BackupOptions) Complete(cmd *cobra.Command, args []string) error {

	opt.command = cmd
	return nil
}

func (opt *BackupOptions) Validate() error {
	args := opt.command.Flags().Args()
	if len(args) > 0 {
		return common.ErrTooManyArguments
	}
	return nil
}

func (opt *BackupOptions) Execute() error {

	sourcePath := hosts.EtcHosts.Path()
	targetPath := opt.output
	if targetPath == "" {
		targetPath = fmt.Sprintf("%s.bak", sourcePath)
	}

	if _, err := os.Stat(targetPath); err == nil {
		if !opt.force {
			return errors.New("backup file already exists")
		}
		os.Remove(targetPath)
	}

	src := hosts.NewSource(sourcePath, common.FileSystem(opt.command.Context()))
	err := src.Apply(func(path string, fs afero.Fs) error {
		sf, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer sf.Close()

		tf, err := fs.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			return err
		}
		defer tf.Close()

		_, err = io.Copy(tf, sf)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
