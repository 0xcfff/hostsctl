package alias

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestAliasDeleteCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "delete - by ip",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/delete__by_ip__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete - by alias",
			Args: cmdtest.ITArgs{
				Args:       []string{"users.example.com"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/delete__by_alias__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete from block - ip by block name",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51", "-b", "pet-prj2"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/delete_from_block__ip_by_block_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete from block - ip by block id",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.51", "-b", "4"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/delete_from_block__ip_by_block_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete many - by alias",
			Args: cmdtest.ITArgs{
				Args:       []string{"awards.example.com"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/delete/delete_many__by_alias___result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestAliasDeleteCommand", func() *cobra.Command { return NewCmdAliasDelete() })
}
