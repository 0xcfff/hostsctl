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
				Args:       []string{"15"},
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/format/format__empty__result.txt",
			},
			Want: true,
		},
		{
			Name: "format - non empty",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "testdata/format/format__non_empty__result.txt",
			},
			Want: true,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestDatabaseFormatCommand", func() *cobra.Command { return NewCmdDatabaseFormat() })
}
