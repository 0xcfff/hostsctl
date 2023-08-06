package alias

import (
	"testing"

	"github.com/0xcfff/hostsctl/testtools"
	"github.com/spf13/cobra"
)

func TestAliasAddCommand(t *testing.T) {
	tests := []testtools.ITTest{
		{
			Name: "empty - args",
			Args: testtools.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/add/empty_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "one line - args",
			Args: testtools.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/one-ip_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "one line - stdin",
			Args: testtools.ITArgs{
				Args:       []string{},
				Stdin:      "127.0.0.1 my.domain.test",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/one-ip_one-alias_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "one line - args + comment",
			Args: testtools.ITArgs{
				Args:       []string{"127.0.0.1", "my.domain.test", "-c", "My custom service domain"},
				Stdin:      "127.0.0.1 my.domain.test",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/one-ip_one-alias-and-comment_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "one line - stdin + comment",
			Args: testtools.ITArgs{
				Args:       []string{},
				Stdin:      "127.0.0.1 my.domain.test # My custom service domain",
				InputFile:  "testdata/one-ip.txt",
				OutputFile: "testdata/add/one-ip_one-alias-and-comment_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// add to specific block
		{
			Name: "add - 3rd block + no block specified",
			Args: testtools.ITArgs{
				Args:       []string{"192.168.100.100", "local-service"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/four-blocks_one-alias_no-block_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - 3rd block + by id",
			Args: testtools.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "3"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/four-blocks_one-alias_3rd-block_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - 3rd block + by name",
			Args: testtools.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "pet-prj1"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/four-blocks_one-alias_3rd-block_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - 5th block + by id",
			Args: testtools.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "5", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/four-blocks_one-alias-by-id_5th-block_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - 5th block + by name",
			Args: testtools.ITArgs{
				Args:       []string{"192.168.100.100", "local-service", "-b", "prj-pet009", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/four-blocks_one-alias-by-name_5th-block_result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},

		// errors cases
		{
			Name: "one line - args + missing block",
			Args: testtools.ITArgs{
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

	testtools.RunIntergationTests(t, tests, "TestAliasAddCommand", func() *cobra.Command { return NewCmdAliasAdd() })
}
