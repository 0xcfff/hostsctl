package block

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestBlockClearCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		// delete empty block
		{
			Name: "clear - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/clear/clear__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj1"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/clear/clear__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear system - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"1", "-f"},
				Stdin:      "",
				InputFile:  "testdata/system-blocks.txt",
				OutputFile: "testdata/clear/clear_system__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear system - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"sys1", "-f"},
				Stdin:      "",
				InputFile:  "testdata/system-blocks.txt",
				OutputFile: "testdata/clear/clear_system__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear not annotated - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"2"},
				Stdin:      "",
				InputFile:  "testdata/three-na-blocks.txt",
				OutputFile: "testdata/clear/clear_not_annotated__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear not annotated - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"2"},
				Stdin:      "",
				InputFile:  "testdata/three-na-blocks.txt",
				OutputFile: "testdata/clear/clear_not_annotated__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "clear error - system",
			Args: cmdtest.ITArgs{
				Args:       []string{"1"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "the block has system aliases",
			},
			Want: false,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestBlockClearCommand", func() *cobra.Command { return NewCmdBlockClear() })
}
