package database

import (
	"fmt"
	"testing"

	"github.com/0xcfff/hostsctl/commands/cmdtest"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/cobra"
)

func TestDatabaseLocationCommand(t *testing.T) {
	tests := []cmdtest.ITTest{
		{
			Name: "location",
			Args: cmdtest.ITArgs{
				Args:      []string{},
				InputFile: "testdata/empty.txt",
				Stdout:    fmt.Sprintf("%s\n", hosts.EtcHosts.Path()),
			},
			Want: true,
		},
		{
			Name: "location error - too many arguments",
			Args: cmdtest.ITArgs{
				Args:      []string{"15"},
				InputFile: "testdata/empty.txt",
				ErrorText: "too many arguments",
			},
			Want: false,
		},
	}
	cmdtest.RunIntergationTests(t, tests, "TestDatabaseLocationCommand", func() *cobra.Command { return NewCmdDatabaseLocation() })
}
