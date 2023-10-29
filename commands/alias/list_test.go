package alias

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestAliasListCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		// default
		{
			Name: "list - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list__empty__output.txt",
			},
			Want: true,
		},
		{
			Name: "list - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/list__two_sys_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// short
		{
			Name: "list short - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "short"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_short__empty__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list short - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "short"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_short__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list short - two sys blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "short"},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/list_short__two_sys_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// wide
		{
			Name: "list wide - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "wide"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_wide__empty__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list wide - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "wide"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_wide__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list wide - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "wide"},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/list_wide__two_sys_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// plain
		{
			Name: "list plain - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "plain"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_plain__empty__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list plain - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "plain"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_plain__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list plain - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "plain"},
				Stdin:      "",
				InputFile:  "testdata/two-sys-blocks.txt",
				StdoutFile: "testdata/list/list_plain__two_sys_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// json set of tests
		{
			Name: "list json - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "json"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_json__empty__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list json - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "json"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_json__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// yaml set of tests
		{
			Name: "list yaml - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "yaml"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				StdoutFile: "testdata/list/list_yaml__empty__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "list yaml - one line",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "yaml"},
				Stdin:      "",
				InputFile:  "testdata/one-ip.txt",
				StdoutFile: "testdata/list/list_yaml__one_ip__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// arrangement cases
		{
			Name: "arange raw - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-a", "raw"},
				Stdin:      "",
				InputFile:  "testdata/two-mixed-blocks.txt",
				StdoutFile: "testdata/list/arange_raw__two_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "arange group - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-a", "group"},
				Stdin:      "",
				InputFile:  "testdata/two-mixed-blocks.txt",
				StdoutFile: "testdata/list/arange_group__two_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "arange ungroup - two blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-a", "ungroup"},
				Stdin:      "",
				InputFile:  "testdata/two-mixed-blocks.txt",
				StdoutFile: "testdata/list/arange_ungroup__two_blocks__output.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},

		// wrong arguments check
		{
			Name: "error - wrong format",
			Args: cmdtest.ITArgs{
				Args:       []string{"-o", "not_supported"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "Error: value not_supported is not support; not supported output format",
			},
			Want: false,
		},
		{
			Name: "error - wrong grouping",
			Args: cmdtest.ITArgs{
				Args:       []string{"-a", "not_supported"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "Error: value not_supported is not support; wrong argument value",
			},
			Want: false,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestAliasListCommand", func() *cobra.Command { return NewCmdAliasList() })
}
