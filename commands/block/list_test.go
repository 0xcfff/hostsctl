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
			Name: "default empty",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/default_empty.txt",
			},
			Want: true,
		},
		{
			Name: "default one line",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/default_one-ip.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "default two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/default_two-sys-blocks.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestBlockListCommand", func() *cobra.Command { return NewCmdBlockList() })
}
