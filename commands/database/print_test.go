package database

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestDatabasePrintCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "print - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/print/print__empty__output.txt",
			},
			Want: true,
		},
		{
			Name: "print - non empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				InputFile:  "testdata/six-blocks.txt",
				StdoutFile: "testdata/print/print__non_empty__output.txt",
			},
			Want: true,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestDatabasePrintCommand", func() *cobra.Command { return NewCmdDatabasePrint() })
}
