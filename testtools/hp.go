package testtools

import (
	"fmt"
	"os"
	"os/exec"
)

// Runs binaries of the current process as sub-prpocess passing extra parameters
// indicating that it is the sub process.
// This functionality is necessary for being able to test functions
// which call exit() internally
func RunHelperProcess(suite string, test string, s ...string) *exec.Cmd {
	cs := []string{fmt.Sprintf("-test.run=%s", suite), "--"}
	cs = append(cs, s...)
	env := []string{
		"GO_TEST_HELPER_PROCESS=1",
		fmt.Sprintf("GO_TEST_SUITE_NAME=%s", suite),
		fmt.Sprintf("GO_TEST_TEST_NAME=%s", test),
	}
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}
