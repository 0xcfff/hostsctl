package alias

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestAliasAddCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "add args - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/add/add_args__empty__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add args - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/add_args__one_ip__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add args - one line + comment",
			Args: cmdtest.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test", "-c", "My custom service domain"},
				Stdin:      "127.0.0.1 my.domain.test",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/add_args__one_ip_and_comment__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add stdin - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "127.0.0.1 my.domain.test",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/add_stdin__one_ip__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add stdin - one line + comment",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "127.0.0.1 my.domain.test # My custom service domain",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/add_stdin__one_ip_and_comment__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// add to specific block
		{
			Name: "add to block - no block specified",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.100", "local-service"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_to_block__no_block_specified__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add to block - #3 by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "3"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_to_block__nr3_by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add to block - #3 by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "pet-prj1"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_to_block__nr3_by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add to block - #5 by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "5", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_to_block__nr5_by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add to block - #5 by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "prj-pet009", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_to_block__nr5_by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},

		// errors cases
		{
			Name: "error - missing block",
			Args: cmdtest.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test", "-b", "local-k8s"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "aliases block 'local-k8s' was not found",
			},
			Want: false,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestAliasAddCommand", func() *cobra.Command { return NewCmdAliasAdd() })
}
