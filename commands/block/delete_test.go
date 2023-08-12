package block

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestBlockDeleteCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		// delete empty block
		{
			Name: "delete - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"18"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/delete/delete__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj3"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/delete/delete__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// delete non-empty block
		{
			Name: "delete non empty - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15", "--force"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/delete/delete_non_empty__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete non empty - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj1", "--force"},
				Stdin:      "",
				InputFile:  "testdata/five-blocks.txt",
				OutputFile: "testdata/delete/delete_non_empty__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// delete many blocks
		{
			Name: "delete many - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15", "--force"},
				Stdin:      "",
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "testdata/delete/delete_many__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "delete many - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj3", "--force"},
				Stdin:      "",
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "testdata/delete/delete_many__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		// errors
		{
			Name: "delete error - no block",
			Args: cmdtest.ITArgs{
				Args:       []string{"99"},
				Stdin:      "",
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "block not found",
			},
			Want: false,
		},
		{
			Name: "delete error - many blocks by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				Stdin:      "",
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "2 blocks found matching parameters: too many entries found",
			},
			Want: false,
		},
		{
			Name: "delete error - many blocks by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj1"},
				Stdin:      "",
				InputFile:  "testdata/six-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "2 blocks found matching parameters: too many entries found",
			},
			Want: false,
		},
		{
			Name: "delete error - non empty block by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "target block is not empty, 1 entry(es) found in the block: too many entries found",
			},
			Want: false,
		},
		{
			Name: "delete error - non empty block by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"pet-prj1"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "target block is not empty, 1 entry(es) found in the block: too many entries found",
			},
			Want: false,
		},
	}

	cmdtest.RunIntergationTests(t, tests, "TestBlockDeleteCommand", func() *cobra.Command { return NewCmdBlockDelete() })
}
