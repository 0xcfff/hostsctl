package block

import (
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/spf13/cobra"
)

func TestBlockAddCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "add - by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"15"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/add/add__by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"mk8s-local"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/add/add__by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add - full data",
			Args: cmdtest.ITArgs{
				Args:       []string{"-n", "15", "-t", "mk8s-local", "-c", "Local Microk8s cluster for a pet project"},
				Stdin:      "",
				InputFile:  "testdata/empty.txt",
				OutputFile: "testdata/add/add__full_data__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},

		// Force flag
		{
			Name: "add force - force update by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"-n", "15", "-t", "pet-prj100", "-c", "Renewed project", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_force__force_update_by_id__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},
		{
			Name: "add force - force update by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"-n", "150", "-t", "pet-prj1", "-c", "Renewed project", "--force"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "testdata/add/add_force__force_update_by_name__result.txt",
				Stdout:     "",
				ErrorText:  "",
			},
			Want: true,
		},

		// Error cases
		{
			Name: "error - already exists by id",
			Args: cmdtest.ITArgs{
				Args:       []string{"-n", "15"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "entry already exists",
			},
			Want: false,
		},
		{
			Name: "error - already exists by name",
			Args: cmdtest.ITArgs{
				Args:       []string{"-t", "pet-prj2"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "entry already exists",
			},
			Want: false,
		},
		{
			Name: "error - too many blocks",
			Args: cmdtest.ITArgs{
				Args:       []string{"-n", "15", "-t", "pet-prj2"},
				Stdin:      "",
				InputFile:  "testdata/four-blocks.txt",
				OutputFile: "",
				Stdout:     "",
				ErrorText:  "multiple blocks found matching criteria",
			},
			Want: false,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestBlockAddCommand", func() *cobra.Command { return NewCmdBlockAdd() })
}
