package alias

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestAliasDeleteCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "four-blocks - by ip",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/four-blocks_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "four-blocks - by ip - in block name",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51", "-b", "pet-prj2"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/four-blocks_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "four-blocks - by ip - in block id",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51", "-b", "4"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/four-blocks_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "four-blocks - by alias",
			Args: cmdtest.ITArgs{
				Args:       []string{"users.example.com"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/four-blocks_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "four-blocks - by alias from multialias line",
			Args: cmdtest.ITArgs{
				Args:       []string{"awards.example.com"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/four-blocks_one-alias-multiline_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestAliasDeleteCommand", func() *cobra.Command { return NewCmdAliasDelete() })
}
