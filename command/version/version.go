package version

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print tool version",
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := debug.ReadBuildInfo()
			if ok {
				fmt.Println("dnssync Version:", b.Main.Version, "Go Version:", b.GoVersion)
			} else {
				cobra.CheckErr(fmt.Errorf("Can't read dnssync version"))
			}

		},
	}
	return cmd
}