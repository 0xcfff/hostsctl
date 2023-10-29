package database

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestDatabaseFormatCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "format - empty",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/format/format__empty__result.txt",
			},
			Want: true,
		},
		{
			Name: "format - non empty",
			Args: cmdtest.ITArgs{
				Args:       []string{},
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "testdata/format/format__non_empty__result.txt",
			},
			Want: true,
		},
		{
			Name: "format dry run - non empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"--dry-run"},
				InputFile:  "testdata/six-blocks.txt",
				StdoutFile: "testdata/format/format_dry_run__non_empty__output.txt",
			},
			Want: true,
		},
		{
			Name: "format error - too many arguments",
			Args: cmdtest.ITArgs{
				Args:      []string{"15"},
				InputFile: "testdata/six-blocks.txt",
				ErrorText: "too many arguments",
			},
			Want: false,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestDatabaseFormatCommand", func() *cobra.Command { return NewCmdDatabaseFormat() })
}
