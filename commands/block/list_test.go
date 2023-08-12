package block

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestBlockListCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		// default
		{
			Name: "list empty",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_empty__output.txt",
			},
			Want: true,
		},
		{
			Name: "list one block",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list two sys blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/list_two_sys_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestBlockListCommand", func() *cobra.Command { return NewCmdBlockList() })
}
