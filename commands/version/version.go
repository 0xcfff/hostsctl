package version

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

type VersionParams struct {
	Version string
}

func NewCmdVersion(p VersionParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print tool version",
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := debug.ReadBuildInfo()
			if ok {
				fmt.Println(os.Args[0], "Version:", p.Version, "Go Version:", b.GoVersion)
			} else {
				fmt.Println(os.Args[0], "Version:", p.Version)
			}

		},
	}
	return cmd
}
