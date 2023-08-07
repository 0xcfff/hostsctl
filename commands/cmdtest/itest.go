package cmdtest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Command integration test arguments and expected outcomes
type ITArgs struct {
	Args       []string
	Stdin      string
	InputFile  string
	OutputFile string
	Stdout     string
	StdoutFile string
	ErrorText  string
}

// Command test case
type ITTest struct {
	Name string
	Args ITArgs
	Want bool
}

// Runs set of tests against command returned by the command factory function cf.
func RunIntergationTests(t *testing.T, tcs []ITTest, tn string, cf func() *cobra.Command) {
	for _, tt := range tcs {
		inHelperProcess := os.Getenv("GO_TEST_HELPER_PROCESS") == "1"
		if inHelperProcess {
			testName := os.Getenv("GO_TEST_TEST_NAME")
			if !strings.EqualFold(testName, tt.Name) {
				continue
			}
		}
		t.Run(tt.Name, func(t *testing.T) {
			if !tt.Want && !inHelperProcess {
				tstp := runHelperProcess(tn, tt.Name)
				out, _ := tstp.CombinedOutput()
				fmt.Println(string(out))
				assert.NotEqual(t, 0, tstp.ProcessState.ExitCode())
				assert.Contains(t, string(out), tt.Args.ErrorText)
				return
			}

			// arrange
			fs := afero.NewMemMapFs()
			fn := hosts.EtcHosts.Path()
			f, err := fs.Create(fn)
			if err != nil {
				t.Errorf("Can't create %v", fn)
				t.FailNow()
			}
			data, err := os.ReadFile(tt.Args.InputFile)
			if err != nil {
				t.Errorf("Can't read %v", tt.Args.InputFile)
				t.FailNow()
			}
			sdata := string(data)
			f.WriteString(sdata)
			f.Close()

			expectDataSpecified := false
			expectData := bytes.NewBufferString("").Bytes()
			if tt.Args.OutputFile != "" {
				expectData, err = os.ReadFile(tt.Args.OutputFile)
				if err != nil {
					t.Errorf("Can't read %v", tt.Args.OutputFile)
					t.FailNow()
				}
				expectDataSpecified = true
			}
			expectRes := string(expectData)
			expectOut := tt.Args.Stdout
			if tt.Args.StdoutFile != "" {
				expectOutBytes, err := os.ReadFile(tt.Args.StdoutFile)
				if err != nil {
					t.Errorf("Can't read %v", tt.Args.OutputFile)
					t.FailNow()
				}
				expectOut = string(expectOutBytes)
			}

			ctx := common.WithCustomFilesystem(context.Background(), fs)
			in := strings.NewReader(tt.Args.Stdin)
			out := &strings.Builder{}

			cmd := cf()
			cmd.SilenceErrors = true
			cmd.SetArgs(tt.Args.Args)
			cmd.SetIn(in)
			cmd.SetOutput(out)
			cmd.SetContext(ctx)

			// act
			c, err := cmd.ExecuteC()

			// assert
			if tt.Want {
				assert.NoError(t, err, "command should succeed")
			} else {
				assert.Error(t, err, "command should fail")
			}

			assert.Same(t, cmd, c)

			s := out.String()
			assert.Equal(t, expectOut, s)

			if expectDataSpecified {
				fr, _ := afero.ReadFile(fs, fn)
				assert.Equal(t, expectRes, string(fr))
			}
		})
	}
}

// Runs binaries of the current process as sub-prpocess passing extra parameters
// indicating that it is the sub process.
// This functionality is necessary for being able to test functions
// which call exit() internally
func runHelperProcess(suite string, test string, s ...string) *exec.Cmd {
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
